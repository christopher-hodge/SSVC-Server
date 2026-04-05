package domain

type AffixType int

const (
	Prefix AffixType = iota
	Suffix
	Either
)

type BaseAffix struct {
	ID       string
	Name     string
	Type     AffixType
	Tags     []string
	BaseMin  int
	BaseMax  int
	Weight   int
	MinLevel int
}

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
	Tier           int
}

func GenerateAffix(base BaseAffix, tier int) AffixDefinition {
	multiplier := 1 + 0.2*float64(tier-1) // +20% per tier

	min := int(float64(base.BaseMin) * multiplier)
	max := int(float64(base.BaseMax) * multiplier)

	weight := int(float64(base.Weight) * (1 - 0.1*float64(tier-1))) // -10% per tier
	if weight < 1 {
		weight = 1
	}

	return AffixDefinition{
		ID:             base.ID,
		Name:           base.Name,
		Type:           base.Type,
		Tags:           base.Tags,
		MinValue:       min,
		MaxValue:       max,
		DisplayedValue: 0,
		Weight:         weight,
		MinLevel:       base.MinLevel,
		Tier:           tier,
	}
}
