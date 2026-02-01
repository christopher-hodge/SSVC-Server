package main

import (
	"log"
	"net/http"

	handlers "SSVC-Server/internal/handlers/http"
)

func main() {

	router := handlers.NewRouter()

	log.Println("Go backend running on :8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}
