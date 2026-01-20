package domain

type Rarity int

const (
	Normal Rarity = iota
	Magic
	Rare
	Unique
)

type Item struct {
	ID        string
	BaseType  string
	Rarity    Rarity
	ItemLevel int

	Prefixes []AffixInstance
	Suffixes []AffixInstance
}
