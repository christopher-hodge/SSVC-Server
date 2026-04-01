package domain

type Rarity int
type MetaModifier int

const (
	Normal Rarity = iota
	Magic
	Rare
	Unique
)

const (
	CannotRollPrefixes MetaModifier = iota
	CannotRollSuffixes
	LockPrefixes
	LockSuffixes
)

type AffixLimits struct {
	MaxImplicits int
	MaxPrefixes  int
	MaxSuffixes  int
}

var AffixLimitsByRarity = map[Rarity]AffixLimits{
	Normal: {
		MaxImplicits: 3,
		MaxPrefixes:  0,
		MaxSuffixes:  0,
	},
	Magic: {
		MaxImplicits: 3,
		MaxPrefixes:  1,
		MaxSuffixes:  1,
	},
	Rare: {
		MaxImplicits: 3,
		MaxPrefixes:  3,
		MaxSuffixes:  3,
	},
	Unique: {
		MaxImplicits: 3,
	},
}

type Item struct {
	ID        string
	BaseType  string
	Rarity    Rarity
	ItemLevel int
	Prefixes  []AffixDefinition
	Suffixes  []AffixDefinition
}

func (i *Item) HasPrefixModifier(id string) bool {
	for _, aff := range i.Prefixes {
		if aff.ID == id {
			return true
		}
	}
	return false
}

func (i *Item) HasSuffixModifier(id string) bool {
	for _, aff := range i.Suffixes {
		if aff.ID == id {
			return true
		}
	}
	return false
}
