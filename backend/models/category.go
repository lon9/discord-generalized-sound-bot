package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lon9/discord-generalized-sound-bot/backend/database"
)

// Category is category of sounds
type Category struct {
	gorm.Model
	Name   string `json:"name"`
	Sounds Sounds `json:"sounds"`
}

// Create creates Category
func (c *Category) Create() (err error) {
	db := database.GetDB()
	return db.Create(c).Error
}

// FindByID finds category by id
func (c *Category) FindByID(id uint) (err error) {
	db := database.GetDB()
	if err = db.First(c, id).Error; err != nil {
		return
	}
	return db.Model(c).Related(&c.Sounds).Error
}

// FindByName finds Category by Name
func (c *Category) FindByName(name string) (err error) {
	db := database.GetDB()
	if err = db.Where("name = ?", name).First(c).Error; err != nil {
		return
	}
	return db.Model(c).Related(&c.Sounds).Error
}

// Categories is slice of Category
type Categories []Category

// FindAll returns all categories
func (cs *Categories) FindAll() (err error) {
	db := database.GetDB()
	return db.Find(cs).Error
}

// SearchByName searches Categories by name
func (cs *Categories) SearchByName(name string) (err error) {
	db := database.GetDB()
	query := fmt.Sprintf("%%%s%%", name)
	return db.Where("name LIKE ?", query).Find(cs).Error
}
