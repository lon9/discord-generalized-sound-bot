package forms

import (
	"database/sql"
	"os"
	"testing"

	"github.com/lon9/discord-generalized-sound-bot/backend/config"
	"github.com/lon9/discord-generalized-sound-bot/backend/database"
	"github.com/lon9/discord-generalized-sound-bot/backend/models"
	testfixtures "gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

func TestMain(m *testing.M) {
	config.Init("unit_test")
	database.Init(true, &models.Sound{}, &models.Category{})
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
