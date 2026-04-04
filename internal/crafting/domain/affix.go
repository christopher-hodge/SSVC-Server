package domain

type AffixType int

const (
	Prefix AffixType = iota
	Suffix
	Either
)

type AffixDefinition struct {
	ID             string
	Name           string
	Type           AffixType
	Tags           []string
	MinValue       int
	MaxValue       int
	DisplayedValue int
	Weight         int
	MinLevel       int
}

type AffixInstance struct {
	DefID string
	Value int
	Tags  []string
}
