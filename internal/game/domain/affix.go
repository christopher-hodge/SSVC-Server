package domain

import (
	"SSVC-Server/internal/random"
)

type AffixType int

const (
	Prefix AffixType = iota
	Suffix
	Either
)

type BaseAffix struct {
	ID         string
	Name       string
	Type       AffixType
	Tags       []string
	ValueCount int
	BaseMin    int
	BaseMax    int
	Weight     int
	MinLevel   int
}

type AffixDefinition struct {
	ID              string
	Name            string
	Type            AffixType
	Tags            []string
	MinValue        int
	MaxValue        int
	DisplayedValues []int
	Weight          int
	MinLevel        int
	Tier            int
}

func GenerateAffix(base BaseAffix, tier int, rng random.RNGer) AffixDefinition {
	multiplier := 1 + 0.2*float64(tier-1) // +20% per tier

	min := int(float64(base.BaseMin) * multiplier)
	max := int(float64(base.BaseMax) * multiplier)

	weight := int(float64(base.Weight) * (1 - 0.1*float64(tier-1))) // -10% per tier
	if weight < 1 {
		weight = 1
	}

	values := rollValues(base.ValueCount, min, max, rng)

	return AffixDefinition{
		ID:              base.ID,
		Name:            base.Name,
		Type:            base.Type,
		Tags:            base.Tags,
		MinValue:        min,
		MaxValue:        max,
		DisplayedValues: values,
		Weight:          weight,
		MinLevel:        base.MinLevel,
		Tier:            tier,
	}
}

func rollValues(count, min, max int, rng random.RNGer) []int {
	switch count {
	case 1:
		return []int{rng.Intn(max-min+1) + min}

	case 2:
		low := rng.Intn(max-min+1) + min
		high := rng.Intn(max-low+1) + low // ensures high ≥ low
		return []int{low, high}

	default:
		return nil
	}
}
