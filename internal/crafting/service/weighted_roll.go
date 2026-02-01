package service

import (
	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/random"
	"time"
)

func weightedRoll(
	pool []domain.AffixDefinition,
) domain.AffixDefinition {

	rng := random.New(time.Now().UnixNano())

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
