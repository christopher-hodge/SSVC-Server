package api

import (
	"encoding/json"
	"net/http"

	"backend/internal/match"
	ws "backend/internal/websocket"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		PlayerID string `json:"playerId"`
	}

	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	json.NewEncoder(w).Encode(map[string]string{
		"token": "mock-token-" + req.PlayerID,
	})
}

func Matchmake(w http.ResponseWriter, r *http.Request) {
	server := match.FindServer()

	json.NewEncoder(w).Encode(server)
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	ws.ServeWebSocket(w, r)
}
