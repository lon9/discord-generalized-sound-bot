package controllers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jonas747/dca"
	filetype "gopkg.in/h2non/filetype.v1"

	"github.com/gin-gonic/gin"
	"github.com/lon9/discord-generalized-sound-bot/backend/config"
	"github.com/lon9/discord-generalized-sound-bot/backend/forms"
	"github.com/lon9/discord-generalized-sound-bot/backend/models"
)

// SoundsController is controller of Sounds
type SoundsController struct{}

// Index returns sounds
func (sc *SoundsController) Index(c *gin.Context) {
	var sounds models.Sounds
	if query := c.Query("query"); query != "" {
		if err := sounds.SearchByName(query); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"result": sounds,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": http.StatusText(http.StatusBadRequest),
	})
}

// Create creates Sound
func (sc *SoundsController) Create(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}
	name := c.PostForm("name")
	categoryName := c.PostForm("categoryName")
	if name == "" || categoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	// Check mime type
	kind, err := filetype.Match(b)
	if err != nil || (kind.Extension != "wav" && kind.Extension != "mp3" && kind.Extension != "ogg" && kind.Extension != "flac") {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// Insert to database
	form := &forms.SoundForm{
		Name:         name,
		CategoryName: categoryName,
	}
	sound, err := form.Create()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	// Encode to dca
	options := dca.StdEncodeOptions
	options.RawOutput = true
	sess, err := dca.EncodeMem(bytes.NewBuffer(b), options)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	defer sess.Cleanup()

	// Save data.
	conf := config.GetConfig()
	dataPath := conf.GetString("data.path")
	savePath := filepath.Join(dataPath, sound.Category.Name, sound.Name+".dca")
	if _, err := os.Stat(filepath.Dir(savePath)); err != nil {
		if err = os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	}
	out, err := os.Create(savePath)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}
	defer out.Close()

	io.Copy(out, sess)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"result": sound,
	})
}
