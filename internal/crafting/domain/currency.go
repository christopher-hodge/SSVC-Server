package domain

import (
	"SSVC-Server/internal/random"
)

type CraftingContext struct {
	Item *Item
	RNG  random.RNGer
}

type Currency interface {
	Apply(ctx *CraftingContext) error
}
