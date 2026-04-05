package service

import (
	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/random"
)

func weightedRoll(
	rng random.RNGer,
	pool []domain.AffixDefinition,
) domain.AffixDefinition {

	totalWeight := 0
	for _, def := range pool {
		totalWeight += def.Weight
	}

	roll := rng.Intn(totalWeight)

	running := 0
	for _, def := range pool {
		running += def.Weight
		if roll < running {
			return def
		}
	}

	return pool[len(pool)-1]
}
