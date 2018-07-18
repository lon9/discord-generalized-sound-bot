package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lon9/discord-generalized-voice-bot/backend/database"
)

// Sound is model of sounds
type Sound struct {
	gorm.Model
	Name       string    `json:"name" gorm:"index"`
	Path       string    `json:"path"`
	CategoryID uint      `json:"categoryId"`
	Category   *Category `json:"category"`
}

// Create creates Sound
func (s *Sound) Create() (err error) {
	db := database.GetDB()
	return db.Create(s).Error
}

// FindByName finds sound by Name
func (s *Sound) FindByName(name string) (err error) {
	db := database.GetDB()
	return db.Where("name = ?", name).First(s).Error
}

// Sounds is slice of Sound
type Sounds []Sound

// FindByCategoryName finds sounds by Category.Name
func (ss *Sounds) FindByCategoryName(categoryName string) (err error) {
	db := database.GetDB()
	var category Category
	if err = db.Where("name = ?", categoryName).First(&category).Error; err != nil {
		return
	}
	return db.Model(&category).Related(ss).Error
}

// SearchByName searches sounds by name
func (ss *Sounds) SearchByName(name string) (err error) {
	db := database.GetDB()
	query := fmt.Sprintf("%%%s%%", name)
	return db.Where("name LIKE ?", query).Find(ss).Error
}
