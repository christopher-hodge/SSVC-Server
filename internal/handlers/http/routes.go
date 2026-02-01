package handlers

import (
	"SSVC-Server/internal/handlers/websocket"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", Health).Methods("GET")
	r.HandleFunc("/matchmake", Matchmake).Methods("POST")
	r.HandleFunc("/ws", websocket.ServeWebSocket)

	return r
}
