package models

import (
	"database/sql"
	"os"
	"testing"

	"gopkg.in/testfixtures.v2"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lon9/discord-generalized-voice-bot/backend/config"
	"github.com/lon9/discord-generalized-voice-bot/backend/database"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	config.Init("unit_test")
	database.Init(true, &Sound{}, &Category{})
	dbURL := config.GetConfig().GetString("db.url")
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		panic(err)
	}

	fixtures, err = testfixtures.NewFolder(db, &testfixtures.SQLite{}, "../fixtures")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}
