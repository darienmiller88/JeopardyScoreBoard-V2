package repositories

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	
	db, err = sqlx.Connect("postgres", os.Getenv("TEST_DATABASE_URL"))

    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    // Optional but recommended:
    // run migrations here

    code := m.Run()

    _ = db.Close()
    os.Exit(code)
}
