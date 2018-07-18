package models

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestCreateSound(t *testing.T) {
	prepareTestDatabase()
	sound := &Sound{
		Name:       "test",
		Path:       "/test/sound1.dca",
		CategoryID: 1,
	}
	if err := sound.Create(); err != nil {
		t.Error(err)
	}

	if sound.ID == 0 {
		t.Error("Sound.ID should not to be 0")
	}

	if sound.Name != "test" {
		t.Errorf("Sound.Name should be test:%s", sound.Name)
	}

	if sound.CategoryID != 1 {
		t.Errorf("Sound.CategoryID should be 1:%d", sound.CategoryID)
	}
}

func TestFindSoundByName(t *testing.T) {
	prepareTestDatabase()
	var sound Sound
	if err := sound.FindByName("Test Sound 1"); err != nil {
		t.Error(err)
	}

	if sound.ID != 1 {
		t.Errorf("Sound.ID should be 1:%d", sound.ID)
	}

	if sound.Name != "Test Sound 1" {
		t.Errorf("Sound.Name should be Test Sound 1:%s", sound.Name)
	}

	if sound.Path != "data/test/sound1.dca" {
		t.Errorf("Sound.Path should be data/test/sound1.dca:%s", sound.Path)
	}

	if sound.CategoryID != 1 {
		t.Errorf("Sound.CategoryID should be 1:%d", sound.CategoryID)
	}

	if err := sound.FindByName("Non-exist sound name"); err == nil {
		t.Error("If it receive non-exist sound name, should return error")
	}
}

func TestFindByCategoryName(t *testing.T) {
	prepareTestDatabase()
	var sounds Sounds
	if err := sounds.FindByCategoryName("Test Category 1"); err != nil {
		t.Error(err)
	}

	if len(sounds) != 2 {
		t.Errorf("Sounds length should be 2:%d", len(sounds))
	}

	if sounds[0].ID != 1 {
		t.Errorf("Sounds[0].ID sould be 1:%d", sounds[0].ID)
	}

	if sounds[0].Name != "Test Sound 1" {
		t.Errorf("Sounds[0].Name sould be Test Sound 1:%s", sounds[0].Name)
	}

	if sounds[0].Path != "data/test/sound1.dca" {
		t.Errorf("Sounds[0].Path sould be data/test/sound1.dca:%s", sounds[0].Path)
	}

	if sounds[0].CategoryID != 1 {
		t.Errorf("Sounds[0].CategoryID should be 1:%d", sounds[0].CategoryID)
	}

	if sounds[1].ID != 2 {
		t.Errorf("Sounds[1].ID sould be 2:%d", sounds[0].ID)
	}

	if sounds[1].Name != "Test Sound 2" {
		t.Errorf("Sounds[1].Name sould be Test Sound 2:%s", sounds[0].Name)
	}

	if sounds[1].Path != "data/test/sound2.dca" {
		t.Errorf("Sounds[1].Path sould be data/test/sound2.dca:%s", sounds[0].Path)
	}

	if sounds[1].CategoryID != 1 {
		t.Errorf("Sounds[1].CategoryID should be 1:%d", sounds[0].CategoryID)
	}

	if err := sounds.FindByCategoryName("Non-exist category"); err == nil {
		t.Error("If it receive non-exist category name, should return error")
	}
}

func TestSearchSoundByName(t *testing.T) {
	prepareTestDatabase()
	var sounds Sounds
	if err := sounds.SearchByName("Test Sound"); err != nil {
		t.Error(err)
	}

	if len(sounds) != 2 {
		t.Errorf("Sounds length should be 2:%d", len(sounds))
	}

	if sounds[0].ID != 1 {
		t.Errorf("Sounds[0].ID sould be 1:%d", sounds[0].ID)
	}

	if sounds[0].Name != "Test Sound 1" {
		t.Errorf("Sounds[0].Name sould be Test Sound 1:%s", sounds[0].Name)
	}

	if sounds[0].Path != "data/test/sound1.dca" {
		t.Errorf("Sounds[0].Path sould be data/test/sound1.dca:%s", sounds[0].Path)
	}

	if sounds[0].CategoryID != 1 {
		t.Errorf("Sounds[0].CategoryID should be 1:%d", sounds[0].CategoryID)
	}

	if sounds[1].ID != 2 {
		t.Errorf("Sounds[1].ID sould be 2:%d", sounds[0].ID)
	}

	if sounds[1].Name != "Test Sound 2" {
		t.Errorf("Sounds[1].Name sould be Test Sound 2:%s", sounds[0].Name)
	}

	if sounds[1].Path != "data/test/sound2.dca" {
		t.Errorf("Sounds[1].Path sould be data/test/sound2.dca:%s", sounds[0].Path)
	}

	if sounds[1].CategoryID != 1 {
		t.Errorf("Sounds[1].CategoryID should be 1:%d", sounds[0].CategoryID)
	}

	if err := sounds.SearchByName("Non-exist sound name"); err != nil {
		t.Error(err)
	}

	if len(sounds) != 0 {
		t.Errorf("If it send non-exist sound name, should receive 0 size slice:%d", len(sounds))
	}
}
