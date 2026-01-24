package domain

type CurrencyType string

const (
	CurrencyImbuementCatalyst      CurrencyType = "imbuement_catalyst"
	CurrencyReconstructionCatalyst CurrencyType = "reconstruction_catalyst"
)

type CraftingContext struct {
	Item        *Item
	AffixRoller AffixRoller
}

type Currency interface {
	Apply(ctx *CraftingContext) error
}
