package domain

type AffixType int

const (
	Prefix AffixType = iota
	Suffix
	Both
	OnlyPrefixes
	OnlySuffixes
	All
)

type AffixDefinition struct {
	ID       string
	Name     string
	Type     AffixType
	Tags     []string
	MinValue int
	MaxValue int
	Weight   int
	MinLevel int
}

type AffixInstance struct {
	DefID string
	Value int
}
