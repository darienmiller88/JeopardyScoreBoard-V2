package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main(){
	godotenv.Load()

	router := chi.NewRouter()

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello world!"))
	})
	
	fmt.Println("Listening on Port:", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}