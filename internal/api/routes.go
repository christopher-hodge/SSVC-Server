package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", Health).Methods("GET")
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/matchmake", Matchmake).Methods("POST")
	r.HandleFunc("/websocket", WebSocket)

	return r
}
