package domain

import (
	"SSVC-Server/internal/random"
)

type CraftingContext struct {
	Item *Item
	RNG  random.RNGerFloat
}

type Currency interface {
	Apply(ctx *CraftingContext) error
}
