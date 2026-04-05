package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"SSVC-Server/internal/game/domain"
	"SSVC-Server/internal/game/service"
	"SSVC-Server/internal/random"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: lock down for prod
	},
}

type WSClient struct {
	Conn *websocket.Conn
	RNG  *random.RNGerFloat
}

type CraftRequest struct {
	Catalyst  string
	AffixType domain.AffixType
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	client := &WSClient{
		Conn: conn,
		RNG:  new(random.RNGerFloat),
	}

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("websocket read failed:", err)
			return
		}

		switch msg.Type {
		case "craft":
			client.handleCraft(msg.Data)

		default:
			_ = conn.WriteJSON(map[string]string{
				"error": "unknown message type",
			})
		}
	}
}

func (c *WSClient) handleCraft(raw json.RawMessage) {

	var req CraftRequest

	if err := json.Unmarshal(raw, &req); err != nil {
		_ = c.Conn.WriteJSON(map[string]string{
			"error": "invalid craft request",
		})
		return
	}

	ctx := &domain.CraftingContext{
		Item: &domain.Item{},
		RNG:  *c.RNG,
	}

	var err error
	switch req.Catalyst {
	case "imbuement":
		err = (&service.ImbuementCatalyst{}).Apply(ctx, req.AffixType, ctx.RNG)

	case "reconstruction":
		err = (&service.ReconstructingCatalyst{}).Apply(ctx, ctx.RNG)

	case "elevating":
		err = (&service.ElevatingCatalyst{}).Apply(ctx, req.AffixType, ctx.RNG)

	case "defiant":
		err = (&service.DefiantCatalyst{}).Apply(ctx, ctx.RNG)

	case "ascendant":
		err = (&service.AscendantCatalyst{}).Apply(ctx, req.AffixType, ctx.RNG)

	case "lustrating":
		err = (&service.LustratingCatalyst{}).Apply(ctx, req.AffixType)

	case "cathartic":
		err = (&service.CatharticCatalyst{}).Apply(ctx, ctx.RNG)

	default:
		err = errors.New("unknown catalyst")
	}

	if err != nil {
		_ = c.Conn.WriteJSON(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_ = c.Conn.WriteJSON(map[string]interface{}{
		"type": "craft_result",
		"item": ctx.Item,
	})
}
