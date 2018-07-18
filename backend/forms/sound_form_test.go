package forms

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestCreateSound(t *testing.T) {

	// New Category and new Sound
	prepareTestDatabase()
	soundForm := &SoundForm{
		Name:         "test",
		CategoryName: "test",
	}
	sound, err := soundForm.Create()
	if err != nil {
		t.Error(err)
	}

	if sound == nil {
		t.Error("Sound should not to be nil")
	}

	if sound.ID == 0 {
		t.Errorf("Sound.ID should not to be 0:%d", sound.ID)
	}

	if sound.Name != "test" {
		t.Errorf("Sound.Name should be test: %s", sound.Name)
	}

	if sound.Path != "sounds_dca/test/test.dca" {
		t.Errorf("Sound.Path should be sounds_dca/test/test.dca:%s", sound.Path)
	}

	if sound.CategoryID == 0 {
		t.Errorf("Sound.CategoryID should not to be 0:%d", sound.CategoryID)
	}

	if sound.Category == nil {
		t.Error("Sound.Category sould not to be nil.")
	}

	if sound.Category.ID == 0 {
		t.Errorf("Sound.Category.ID should not to be 0:%d", sound.Category.ID)
	}

	if sound.Category.Name != "test" {
		t.Errorf("Sound.Category.Name should be test:%s", sound.Category.Name)
	}

	// New Sound
	prepareTestDatabase()
	soundForm = &SoundForm{
		Name:         "test",
		CategoryName: "Test Category 1",
	}
	sound, err = soundForm.Create()
	if err != nil {
		t.Error(err)
	}

	if sound == nil {
		t.Error("Sound should not to be nil")
	}

	if sound.ID == 0 {
		t.Errorf("Sound.ID should not to be 0:%d", sound.ID)
	}

	if sound.Name != "test" {
		t.Errorf("Sound.Name should be test: %s", sound.Name)
	}

	if sound.Path != "sounds_dca/Test Category 1/test.dca" {
		t.Errorf("Sound.Path should be sounds_dca/Test Category 1/test.dca:%s", sound.Path)
	}

	if sound.CategoryID != 1 {
		t.Errorf("Sound.CategoryID should be 1:%d", sound.CategoryID)
	}

	if sound.Category == nil {
		t.Error("Sound.Category sould not to be nil.")
	}

	if sound.Category.ID != 1 {
		t.Errorf("Sound.Category.ID should be 1:%d", sound.Category.ID)
	}

	if sound.Category.Name != "Test Category 1" {
		t.Errorf("Sound.Category.Name should be Test Category 1:%s", sound.Category.Name)
	}

	// Exist Sound
	prepareTestDatabase()
	soundForm = &SoundForm{
		Name:         "Test Sound 1",
		CategoryName: "Test Category 1",
	}
	sound, err = soundForm.Create()
	if err == nil {
		t.Error("If receive exist-sound-name, should return error")
	}

}
