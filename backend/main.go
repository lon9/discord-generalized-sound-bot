package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lon9/discord-generalized-sound-bot/backend/config"
	"github.com/lon9/discord-generalized-sound-bot/backend/database"
	"github.com/lon9/discord-generalized-sound-bot/backend/models"
	"github.com/lon9/discord-generalized-sound-bot/backend/server"
)

func main() {
	env := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: backend -e {mode}")
		os.Exit(1)
	}
	flag.Parse()

	config.Init(*env)
	database.Init(false, &models.Sound{}, &models.Category{})
	defer database.Close()
	server.Init()
}
