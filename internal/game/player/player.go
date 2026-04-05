package player

import (
	"SSVC-Server/internal/game/domain"
	"SSVC-Server/internal/game/modifier"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Player struct {
	ID        string
	Position  Vector3
	Stats     Stats
	Resources Resources
	State     PlayerState
	Modifiers []modifier.Modifier
	Inventory []domain.Item
	Abilities []Ability
	Cooldowns map[string]float64
}

type Ability struct {
}

type PlayerState struct {
	InCover bool
}
