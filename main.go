package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"JeopardyScoreBoardV2/database"
	"JeopardyScoreBoardV2/controllers"
)

func main(){
	//Load env file immediately at the start of the program
	godotenv.Load()

	//Create new chi router instance to push handlers to.
	router := chi.NewRouter()

	//Middleware stack, keeping it basic for now.
	router.Use(middleware.Logger, middleware.Recoverer)
	
	//Initiate the database connection to mongodb, and defer its disconnection.
	database.Init()
	defer database.DisconnectClient()

	//Initialize the parent controller router, and its children
	index := controllers.Index{}
	index.Init()

	//Afterwards, mount that router onto this one.
	router.Mount("/", index.Router)

	//Finally, listen and serve on the port in the env, which is 8080 on local machine.
	fmt.Println("Listening on Port:", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}