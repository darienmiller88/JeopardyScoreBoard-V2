package repositories

import (
	"JeopardyScoreBoardV2/encryption"
	"encoding/base64"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB
var es *encryption.EncryptionService

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	var err error	

	db, err = sqlx.Connect("postgres", os.Getenv("TEST_DATABASE_URL"))

    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

	key := os.Getenv("ENCRYPTION_KEY")
	keyB64, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		log.Fatal(err)
	}

	es = encryption.NewService(keyB64)
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(
        "file://../test-migrations",
        "postgres", 
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	//Run all of the migrations to recreate the production database
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
	
	code := m.Run()
	migrations.Down()	
	_ = db.Close()
	
	os.Exit(code)
	// defer func ()  {
	// } ()
}

