package service

import (
	"errors"

	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/crafting/tables"
)

func RollAffix(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) (domain.AffixInstance, error) {

	candidates := make([]domain.AffixDefinition, 0)

	for _, def := range tables.AffixPool {
		if def.Type != affixType {
			continue
		}

		if ctx.Item.ItemLevel < def.MinLevel {
			continue
		}

		// future: tag checks, meta-mod blocking, influence, etc.
		candidates = append(candidates, def)
	}

	if len(candidates) == 0 {
		return domain.AffixInstance{}, errors.New("no valid affixes")
	}

	chosen := weightedRoll(ctx.RNG, candidates)

	value := ctx.RNG.Intn(chosen.MaxValue-chosen.MinValue+1) + chosen.MinValue

	return domain.AffixInstance{
		DefID: chosen.ID,
		Value: value,
	}, nil
}
