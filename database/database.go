package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func Init(){
	_db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	if err != nil{
		panic(err)
	}
	
	err = _db.Ping()
	
	if err != nil {
		fmt.Println("db connection fail:", err)
	}else{
		fmt.Println("Connection established! :)")
	}
	
	db = _db
}
	
func GetDB() *sqlx.DB{
	return db
}

func CloseSQLDB(){
	if err := db.Close(); err != nil{
		panic(err)
	}
}