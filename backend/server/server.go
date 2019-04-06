package server

import "github.com/lon9/discord-generalized-sound-bot/backend/config"

// Init initializes server
func Init() error {
	config := config.GetConfig()
	r, err := NewRouter()
	if err != nil {
		return err
	}
	r.Run(":" + config.GetString("server.port"))
	return nil
}
