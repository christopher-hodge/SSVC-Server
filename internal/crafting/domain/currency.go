package domain

type CurrencyType string

const (
	CurrencyTransmutation CurrencyType = "orb_of_transmutation"
	CurrencyAlteration    CurrencyType = "orb_of_alteration"
)

type CraftingContext struct {
	RNG  RandomSource
	Item *Item
}

type Currency interface {
	Apply(ctx *CraftingContext) error
}
