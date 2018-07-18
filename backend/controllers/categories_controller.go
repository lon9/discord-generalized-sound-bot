package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"github.com/lon9/discord-generalized-voice-bot/backend/models"
)

// CategoriesController is controller for Category
type CategoriesController struct{}

// Index returns categories
func (cc CategoriesController) Index(c *gin.Context) {
	var categories models.Categories
	if query := c.Query("query"); query != "" {
		if err := categories.SearchByName(query); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"result": categories,
		})
		return
	}
	if err := categories.FindAll(); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"result": categories,
	})
}

// Show returns category
func (cc CategoriesController) Show(c *gin.Context) {
	var category models.Category
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}
	if err := category.FindByID(uint(id)); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": http.StatusText(http.StatusNotFound),
			})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"result": &category,
	})
}
