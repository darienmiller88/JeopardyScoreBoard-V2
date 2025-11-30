package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"JeopardyScoreBoardV2/database"
)

func main(){
	godotenv.Load()

	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.Recoverer)
	
	database.Init()
	defer database.DisconnectClient()

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello world!"))
	})

	
	fmt.Println("Listening on Port:", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}