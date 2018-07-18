package server

import "github.com/lon9/discord-generalized-voice-bot/backend/config"

// Init initializes server
func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.Run(":" + config.GetString("server.port"))
}
