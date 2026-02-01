package handlers

import (
	"encoding/json"
	"net/http"

	"SSVC-Server/internal/match"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Matchmake(w http.ResponseWriter, r *http.Request) {
	server := match.FindServer()
	json.NewEncoder(w).Encode(server)
}
