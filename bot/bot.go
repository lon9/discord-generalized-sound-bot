package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jonas747/dca"
	"github.com/jonas747/dshardmanager"
	cache "github.com/patrickmn/go-cache"
	yaml "gopkg.in/yaml.v2"
)

// Silence is silense packet
var Silence = []byte{0xF8, 0xFF, 0xFE}

// Bot is instance of bot
type Bot struct {
	dshardmanager.Manager
	db     *gorm.DB
	config *Config
	queues map[string][]*Play
	m      *sync.Mutex
	cache  *cache.Cache
}

// Play is object for playing sound
type Play struct {
	GuildID   string
	ChannelID string
	UserID    string
	Sound     *Sound
}

// NewBot is constructor
func NewBot(config *Config) (bot *Bot, err error) {
	bot = &Bot{
		Manager: *dshardmanager.New("Bot " + config.Token),
		config:  config,
		queues:  make(map[string][]*Play, config.MaxQueueSize),
		m:       new(sync.Mutex),
		cache:   cache.New(15*time.Minute, 1*time.Minute),
	}

	// Open database
	bot.db, err = gorm.Open("sqlite3", config.BDURL)
	if err != nil {
		return
	}

	// Setting Discord sessions
	bot.Name = config.Name
	if config.Env == "development" {

		// In development
		bot.SetNumShards(1)
		bot.db.LogMode(true)
		bot.db.AutoMigrate(&Sound{})
	} else {

		// In production
		recommended, err := bot.GetRecommendedCount()
		if err != nil {
			return nil, err
		}
		if recommended < 2 {
			bot.SetNumShards(5)
		}
		if config.LogChannelID != "" {
			bot.LogChannel = config.LogChannelID
		}
	}
	bot.AddHandler(bot.ready)
	bot.AddHandler(bot.messageCreate)

	bot.cache.OnEvicted(func(key string, data interface{}) {
		log.Printf("Evicted cache:%s\n", key)
	})

	return bot, bot.Start()
}

// Close closes discord session and db connection
func (b *Bot) Close() {
	b.StopAll()
	b.db.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, b.config.BotPlaying)
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!invite "+b.config.Name) {
		s.ChannelMessageSend(m.ChannelID, "https://discordapp.com/oauth2/authorize?client_id="+b.config.ApplicationID+"&scope=bot&permissions=0")
		return
	}

	if strings.HasPrefix(m.Content, "!help "+b.config.Name) {
		s.ChannelMessageSend(m.ChannelID, b.config.BotHelp)
		return
	}

	if strings.HasPrefix(m.Content, b.config.BotPrefix) {

		cmdName := strings.Replace(m.Content, b.config.BotPrefix, "", -1)
		var sound Sound
		if b.db.Where("name = ?", cmdName).First(&sound).RecordNotFound() {
			s.ChannelMessageSend(m.ChannelID, cmdName+b.config.BotNotFound)
			return
		}

		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Println(err)
			return
		}
		guild, err := s.State.Guild(channel.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		b.requestPlay(m.Author, guild, &sound)
	}
}

func (b *Bot) createPlay(user *discordgo.User, guild *discordgo.Guild, sound *Sound) (*Play, error) {
	channel, err := b.getCurrentVoiceChannel(user, guild)
	if err != nil {
		return nil, err
	}
	play := &Play{
		GuildID:   guild.ID,
		ChannelID: channel.ID,
		UserID:    user.ID,
		Sound:     sound,
	}
	return play, err
}

func (b *Bot) getCurrentVoiceChannel(user *discordgo.User, guild *discordgo.Guild) (*discordgo.Channel, error) {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == user.ID {
			channel, err := b.SessionForGuildS(guild.ID).State.Channel(vs.ChannelID)
			if err != nil {
				return nil, err
			}
			return channel, err
		}
	}
	return nil, errors.New("There is not user in voice channel")
}

func (b *Bot) requestPlay(user *discordgo.User, guild *discordgo.Guild, sound *Sound) {
	play, err := b.createPlay(user, guild, sound)
	if err != nil {
		log.Println(err)
		return
	}

	b.m.Lock()
	if queue, ok := b.queues[guild.ID]; ok {
		if len(queue) < b.config.MaxQueueSize {
			b.queues[guild.ID] = append(queue, play)
		}
	} else {
		b.queues[guild.ID] = []*Play{play}
		go b.runPlayer(guild.ID)
	}
	b.m.Unlock()
}

func (b *Bot) runPlayer(guildID string) {
	var lastChannel string
	var vc *discordgo.VoiceConnection
	for {
		b.m.Lock()
		var play *Play

		if queue, ok := b.queues[guildID]; ok && len(queue) > 0 {
			play = queue[0]
			b.queues[guildID] = queue[1:]
		} else {
			break
		}
		b.m.Unlock()

		if lastChannel != play.ChannelID && vc != nil {
			vc.Disconnect()
			vc = nil
		}

		var err error
		vc, err = b.playSound(play, vc)
		if err != nil {
			log.Println(err)
		}
		lastChannel = play.ChannelID
	}
	if vc != nil {
		vc.Disconnect()
	}

	delete(b.queues, guildID)
	b.m.Unlock()
}

func (b *Bot) playSound(play *Play, vc *discordgo.VoiceConnection) (*discordgo.VoiceConnection, error) {

	var err error
	if vc == nil || !vc.Ready {
		vc, err = b.SessionForGuildS(play.GuildID).ChannelVoiceJoin(play.GuildID, play.ChannelID, false, true)
		if err != nil {
			// return nil, err
			// TODO panic because of unstable behavior.
			panic(err)
		}
	}

	var data [][]byte
	if raw, ok := b.cache.Get(play.Sound.Name); ok {

		// Cache hit
		data = raw.([][]byte)
		vc.Speaking(true)
		if err := b.sendSilence(vc, 10); err != nil {
			return nil, err
		}
		for i := range data {
			vc.OpusSend <- data[i]
		}
		if err := b.sendSilence(vc, 5); err != nil {
			return nil, err
		}

		// Update expiration
		if b.cache.ItemCount() < b.config.SoundCacheSize {
			if err := b.cache.Replace(play.Sound.Name, data, cache.DefaultExpiration); err != nil {
				return nil, err
			}
		}
		return vc, nil
	}

	// Load from file
	f, err := os.Open(filepath.Join(b.config.SoundDir, play.Sound.Path))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	vc.Speaking(true)
	if err = b.sendSilence(vc, 10); err != nil {
		return nil, err
	}
	decoder := dca.NewDecoder(f)
	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		data = append(data, frame)
		vc.OpusSend <- frame
	}

	if err = b.sendSilence(vc, 5); err != nil {
		return nil, err
	}

	if b.cache.ItemCount() < b.config.SoundCacheSize {
		b.cache.SetDefault(play.Sound.Name, data)
	}

	return vc, err
}

func (b *Bot) sendSilence(vc *discordgo.VoiceConnection, n int) (err error) {
	for i := 0; i < n; i++ {
		vc.OpusSend <- Silence
	}
	return
}

// Config is configuration of Bot
type Config struct {
	Name           string `yaml:"name"`
	ApplicationID  string `yaml:"applicationId"`
	Token          string `yaml:"botToken"`
	BDURL          string `yaml:"dbUrl"`
	BotPrefix      string `yaml:"botPrefix"`
	BotHelp        string `yaml:"botHelp"`
	BotPlaying     string `yaml:"botPlaying"`
	BotNotFound    string `yaml:"botNotFound"`
	SoundDir       string `yaml:"soundDir"`
	SoundCacheSize int    `yaml:"soundCacheSize"`
	MaxQueueSize   int    `yaml:"maxQueueSize"`
	Env            string `yaml:"env"`
	LogChannelID   string `yaml:"logChannelId"`
}

// NewConfig is constructor
func NewConfig(fname string) (config *Config, err error) {
	b, err := ioutil.ReadFile(fname)
	err = yaml.Unmarshal(b, &config)
	return
}

// Sound is structure of sounds table
type Sound struct {
	ID   int
	Name string `gorm:"idnex"`
	Path string
}
