package repositories

import (
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
	defer func ()  {
		migrations.Down()	
		_ = db.Close()
		
		os.Exit(code)
	} ()
}

