package main

import (
	"log"
	"net/http"

	"SSVC-Server/internal/api"
)

func main() {
	router := api.NewRouter()

	log.Println("Go backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
