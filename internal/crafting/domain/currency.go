package domain

import (
	"SSVC-Server/internal/random"
)

type CurrencyType string

const (
	CurrencyImbuementCatalyst      CurrencyType = "imbuement_catalyst"
	CurrencyReconstructionCatalyst CurrencyType = "reconstruction_catalyst"
)

type CraftingContext struct {
	Item   *Item
	RNG    *random.RNG
	Tables []AffixDefinition
}

type Currency interface {
	Apply(ctx *CraftingContext) error
}
