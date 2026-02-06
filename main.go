package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"JeopardyScoreBoardV2/controllers"
	"JeopardyScoreBoardV2/database"
	"JeopardyScoreBoardV2/encryption"
)

func main() {
	//Load env file immediately at the start of the program
	godotenv.Load()

	//Create new chi router instance to push handlers to.
	router := chi.NewRouter()

	//Middleware stack, keeping it basic for now.
	router.Use(middleware.Logger, middleware.Recoverer)

	//Initiate the database connection to SQL, and defer its disconnection.
	database.Init()
	defer database.CloseSQLDB()

	key := os.Getenv("ENCRYPTION_KEY")
	keyB64, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		log.Fatal(err)
	}

	//Initialize encryption service
	encryptionService := encryption.NewService(keyB64)

	//Initialize the parent controller router, and its children
	index := controllers.Index{}
	index.InitControllers(database.GetDB(), encryptionService)

	//Afterwards, mount that router onto this one.
	router.Mount("/", index.Router)

	//Serve static files along the "/static" route
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	//Finally, listen and serve on the port in the env, which is 8080 on local machine.
	fmt.Println("Listening on Port:", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}
