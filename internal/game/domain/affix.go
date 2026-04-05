package domain

import (
	"SSVC-Server/internal/game/modifier"
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
	ValueCount float64
	BaseMin    float64
	BaseMax    float64
	Weight     int
	MinLevel   int
}

type AffixDefinition struct {
	ID              string
	Name            string
	Type            AffixType
	Tags            []string
	MinValue        float64
	MaxValue        float64
	DisplayedValues []float64
	Modifiers       []modifier.Modifier
	Weight          int
	MinLevel        int
	Tier            int
}

func GenerateAffix(base BaseAffix, tier int, rng random.RNGerFloat) AffixDefinition {
	multiplier := 1 + 0.2*float64(tier-1) // +20% per tier

	min := float64(base.BaseMin) * multiplier
	max := float64(base.BaseMax) * multiplier

	weight := int(float64(base.Weight) * (1 - 0.1*float64(tier-1))) // -10% per tier
	if weight < 1 {
		weight = 1
	}

	values := rollValues(base.ValueCount, min, max, rng)

	affix := AffixDefinition{
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

	affix.Modifiers = modifier.BuildModifiers(base.ID, []float64(values))

	return affix
}

func rollValues(count, min, max float64, rng random.RNGerFloat) []float64 {
	switch count {
	case 1:
		num := max - min + 1
		return []float64{rng.Floatn(num)}

	case 2:

		low := rng.Floatn((max - min + 1) + min)
		high := rng.Floatn((max - low + 1) + low)
		return []float64{low, high}

	default:
		return nil
	}
}
