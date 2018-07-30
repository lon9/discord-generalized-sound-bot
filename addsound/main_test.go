package main

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lon9/discord-generalized-sound-bot/backend/models"
)

var (
	db *gorm.DB
)

func TestAddSounds(t *testing.T) {
	opts := &options{
		SrcDir:     "src",
		DistDir:    "dist",
		PathPrefix: "/data",
	}

	if err := addSounds(opts); err != nil {
		t.Error(err)
	}

	if _, err := os.Stat("dist/category1/sample1.dca"); err != nil {
		t.Error("should be exist dist/category1/sample1.dca")
	}

	if _, err := os.Stat("dist/category1/sample2.dca"); err != nil {
		t.Error("should be exist dist/category1/sample2.dca")
	}

	if _, err := os.Stat("dist/category1/sample3.dca"); err != nil {
		t.Error("should be exist dist/category1/sample3.dca")
	}

	if _, err := os.Stat("dist/category2/sample4.dca"); err != nil {
		t.Error("should be exist dist/category2/sample4.dca")
	}

	if _, err := os.Stat("dist/category1/sample5.dca"); err == nil {
		t.Error("should not be exist dist/category1/sample5.dca")
	}

	var categories models.Categories
	if err := db.Find(&categories).Error; err != nil {
		t.Error(err)
	}

	if len(categories) != 2 {
		t.Errorf("should be length of Categories is 2:%d", len(categories))
	}

	var sounds models.Sounds
	if err := db.Find(&sounds).Error; err != nil {
		t.Error(err)
	}

	if len(sounds) != 4 {
		t.Errorf("sound be length of Sounds is 4:%d", len(sounds))
	}

	if err := addSounds(opts); err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open("sqlite3", "dist/sounds.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.DropTableIfExists("sounds", "categories")
	if err := os.RemoveAll("dist/category1"); err != nil {
		panic(err)
	}
	if err := os.RemoveAll("dist/category2"); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
