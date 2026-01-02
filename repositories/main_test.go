package repositories

import (
	"fmt"
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
		fmt.Println("test database:", os.Getenv("TEST_DATABASE_URL"))
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

	driver, _ := postgres.WithInstance(db.DB, &postgres.Config{})
	migrations, err := migrate.NewWithDatabaseInstance(
        "file://../migrations",
        "postgres", 
		driver,
	)

	if err != nil{
		fmt.Println("migrations err:")
		panic(err)
	}

	//Run all of the migrations to recreate the production database
    migrations.Up()

    code := m.Run()

	//Clean up by removing the tables in the test database.
	migrations.Down()

    _ = db.Close()
    os.Exit(code)
}
