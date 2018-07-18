package forms

import (
	"errors"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/lon9/discord-generalized-voice-bot/backend/config"
	"github.com/lon9/discord-generalized-voice-bot/backend/models"
)

// SoundForm is form to add sound
type SoundForm struct {
	Name         string `json:"name"`
	CategoryName string `json:"categoryName"`
}

// Create creates a sound
func (sf *SoundForm) Create() (ret *models.Sound, err error) {
	var sound models.Sound
	var category models.Category
	err = sound.FindByName(sf.Name)
	if err == nil && sound.ID != 0 {
		return nil, errors.New("Sound name is duplicated")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err = category.FindByName(sf.CategoryName); err != nil {
		if err == gorm.ErrRecordNotFound {

			// Record not found
			category.Name = sf.CategoryName
			ret = &models.Sound{
				Name:     sf.Name,
				Path:     filepath.Join(config.GetConfig().GetString("data.prefix"), sf.CategoryName, sf.Name+".dca"),
				Category: &category,
			}
			err = ret.Create()
			return
		}
		// Unexpected error
		return nil, err
	}

	ret = &models.Sound{
		Name:       sf.Name,
		Path:       filepath.Join(config.GetConfig().GetString("data.prefix"), sf.CategoryName, sf.Name+".dca"),
		CategoryID: category.ID,
	}
	err = ret.Create()
	ret.Category = &category
	return
}
