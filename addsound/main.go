package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jonas747/dca"
	filetype "gopkg.in/h2non/filetype.v1"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lon9/discord-generalized-sound-bot/backend/models"
)

type options struct {
	SrcDir     string
	DistDir    string
	PathPrefix string
}

func parseFlag() (opts *options) {
	opts = new(options)
	flag.StringVar(&opts.SrcDir, "s", "src", "Source directory")
	flag.StringVar(&opts.DistDir, "d", "dist", "Distination directory")
	flag.StringVar(&opts.PathPrefix, "p", "/sounds_dca", "Prefix of the path to save sounds")
	flag.Parse()
	return
}

func main() {
	opts := parseFlag()
	if err := addSounds(opts); err != nil {
		log.Fatal(err)
	}
}

func addSounds(opts *options) (err error) {
	db, err := gorm.Open("sqlite3", filepath.Join(opts.DistDir, "sounds.db"))
	if err != nil {
		return err
	}
	defer db.Close()
	db.AutoMigrate(new(models.Sound), new(models.Category))

	return filepath.Walk(opts.SrcDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if !info.IsDir() {

			// It's file

			soundName := filepath.Base(p[:len(p)-len(filepath.Ext(p))])
			categoryName := filepath.Base(filepath.Dir(p))
			dataPath := filepath.Join(opts.DistDir, categoryName)
			savePath := filepath.Join(dataPath, soundName+".dca")

			var isError bool
			var isMkDir bool
			defer func() {

				// Teerdown
				if isError {
					if isMkDir {
						if err := os.RemoveAll(dataPath); err != nil {
							log.Println(err)
						}
					} else {
						if err := os.Remove(savePath); err != nil {
							log.Println(err)
						}
					}
				}
			}()

			// Check duplicate.
			if err := db.Where("name = ?", soundName).First(new(models.Sound)).Error; err == nil {
				// Already exists
				err = errors.New("Already exists")
				log.Println(err)
				return nil
			} else if !gorm.IsRecordNotFoundError(err) {
				log.Println(err)
				return nil
			}

			// Check file type.
			kind, err := filetype.MatchFile(p)
			if err != nil {
				log.Println(err)
				return nil
			}
			if kind.Extension != "wav" && kind.Extension != "mp3" && kind.Extension != "ogg" && kind.Extension != "flac" {
				err = fmt.Errorf("Invalid format %s", kind.Extension)
				log.Println(err)
				return nil
			}

			// Transaction
			tx := db.Begin()
			var category models.Category
			if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {

					// If there is no category in database, createit
					sound := &models.Sound{
						Name: soundName,
						Path: filepath.Join(opts.PathPrefix, categoryName, soundName+".dca"),
						Category: &models.Category{
							Name: categoryName,
						},
					}
					if err := tx.Create(sound).Error; err != nil {
						log.Println(err)
						tx.Rollback()
						return nil
					}
				} else {

					// Unexpected error
					log.Println(err)
					tx.Rollback()
					return nil
				}
			} else {

				// There is the category
				sound := &models.Sound{
					Name:       soundName,
					Path:       filepath.Join(opts.PathPrefix, categoryName, soundName+".dca"),
					CategoryID: category.ID,
				}

				if err := tx.Create(sound).Error; err != nil {
					log.Println(err)
					tx.Rollback()
					return nil
				}
			}
			// Commit.
			tx.Commit()

			// Save file
			f, err := os.Open(p)
			if err != nil {
				log.Println(err)
				return nil
			}
			defer f.Close()

			encodeOpts := dca.StdEncodeOptions
			encodeOpts.RawOutput = true
			sess, err := dca.EncodeMem(f, encodeOpts)
			if err != nil {
				log.Println(err)
				return nil
			}
			defer sess.Cleanup()

			if _, err := os.Stat(dataPath); err != nil {
				if err = os.MkdirAll(dataPath, os.ModePerm); err != nil {
					isError = true
					log.Println(err)
					return nil
				}
				isMkDir = true
			}
			out, err := os.Create(savePath)
			if err != nil {
				isError = true
				log.Println(err)
				return nil
			}
			defer out.Close()
			if _, err := io.Copy(out, sess); err != nil {
				isError = true
				return nil
			}
		}
		return nil
	})
}
