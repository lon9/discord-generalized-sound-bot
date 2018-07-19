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
	cache "github.com/patrickmn/go-cache"
	yaml "gopkg.in/yaml.v2"
)

// Bot is instance of bot
type Bot struct {
	dg     *discordgo.Session
	db     *gorm.DB
	config *Config
	queues *sync.Map
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
		config: config,
		queues: &sync.Map{},
		cache:  cache.New(15*time.Minute, 1*time.Minute),
	}
	bot.db, err = gorm.Open("sqlite3", config.BDURL)
	if err != nil {
		return
	}
	bot.dg, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		return
	}
	bot.dg.AddHandler(bot.ready)
	bot.dg.AddHandler(bot.guildCreate)
	bot.dg.AddHandler(bot.messageCreate)

	bot.cache.OnEvicted(func(key string, data interface{}) {
		log.Printf("Evicted cache:%s\n", key)
	})

	return bot, bot.dg.Open()
}

// Close closes discord session and db connection
func (b *Bot) Close() {
	b.dg.Close()
	b.db.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, b.config.BotPlaying)
}

func (b *Bot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			s.ChannelMessageSend(channel.ID, b.config.BotHello)
			return
		}
	}
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
		s.ChannelMessageSend(m.ChannelID, b.config.BotHello)
		return
	}

	if strings.HasPrefix(m.Content, b.config.BotPrefix) {

		cmdName := strings.Replace(m.Content, b.config.BotPrefix, "", -1)
		var sound Sound
		if b.db.Where("name = ?", cmdName).First(&sound).RecordNotFound() {
			s.ChannelMessageSend(m.ChannelID, cmdName+b.config.BotNotFound)
			return
		}

		channel, err := b.dg.State.Channel(m.ChannelID)
		if err != nil {
			log.Println(err)
			return
		}
		guild, err := b.dg.State.Guild(channel.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		go b.enqueuePlay(m.Author, guild, &sound)
	}
}

func (b *Bot) voiceLoop(guildID string, ch chan *Play) {
	var vc *discordgo.VoiceConnection
	var err error
	for {
		select {
		case play := <-ch:
			if vc == nil {
				vc, err = b.dg.ChannelVoiceJoin(play.GuildID, play.ChannelID, false, false)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			if vc.ChannelID != play.ChannelID {
				if err = vc.ChangeChannel(play.ChannelID, false, false); err != nil {
					log.Println(err)
					continue
				}
			}

			// Wait for join
			time.Sleep(time.Millisecond * 125)

			if err := b.playSound(play, vc); err != nil {
				log.Println(err)
				continue
			}
		case <-time.After(1 * time.Second):

			// Disconnect from channel after 1 second
			b.queues.Delete(guildID)
			close(ch)
			vc.Disconnect()
			vc.Close()
			return
		}
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
			channel, err := b.dg.State.Channel(vs.ChannelID)
			if err != nil {
				return nil, err
			}
			return channel, err
		}
	}
	return nil, errors.New("There is not user in voice channel")
}

func (b *Bot) enqueuePlay(user *discordgo.User, guild *discordgo.Guild, sound *Sound) {
	play, err := b.createPlay(user, guild, sound)
	if err != nil {
		log.Println(err)
		return
	}

	if v, ok := b.queues.Load(guild.ID); ok {
		c := v.(chan *Play)
		if len(c) < b.config.MaxQueueSize {
			c <- play
		}
	} else {
		c := make(chan *Play, b.config.MaxQueueSize)
		go b.voiceLoop(guild.ID, c)
		b.queues.Store(guild.ID, c)
		c <- play
	}
}

func (b *Bot) playSound(play *Play, vc *discordgo.VoiceConnection) (err error) {

	var data [][]byte
	if raw, ok := b.cache.Get(play.Sound.Name); ok {

		// Cache hit
		data = raw.([][]byte)
		vc.Speaking(true)
		time.Sleep(time.Millisecond * 125)
		for i := range data {
			vc.OpusSend <- data[i]
		}
		vc.Speaking(false)
		time.Sleep(time.Millisecond * 125)

		// Update expiration
		if b.cache.ItemCount() < b.config.SoundCacheSize {
			if err = b.cache.Replace(play.Sound.Name, data, cache.DefaultExpiration); err != nil {
				return
			}
		}
		return
	}

	// Load from file
	f, err := os.Open(filepath.Join(b.config.SoundDir, play.Sound.Path))
	if err != nil {
		return
	}

	vc.Speaking(true)
	time.Sleep(time.Millisecond * 125)
	decoder := dca.NewDecoder(f)
	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		data = append(data, frame)
		select {
		case vc.OpusSend <- frame:
		case <-time.After(time.Second):
			return err
		}
	}

	vc.Speaking(false)
	time.Sleep(time.Millisecond * 125)

	if b.cache.ItemCount() < b.config.SoundCacheSize {
		b.cache.SetDefault(play.Sound.Name, data)
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
	BotHello       string `yaml:"botHello"`
	BotPlaying     string `yaml:"botPlaying"`
	BotNotFound    string `yaml:"botNotFound"`
	SoundDir       string `yaml:"soundDir"`
	SoundCacheSize int    `yaml:"soundCacheSize"`
	MaxQueueSize   int    `yaml:"maxQueueSize"`
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
	Name string
	Path string
}
