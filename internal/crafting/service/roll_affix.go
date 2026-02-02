package service

import (
	"errors"

	"SSVC-Server/internal/crafting/domain"
)

func RollAffix(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) (domain.AffixInstance, error) {

	validAffixes := make([]domain.AffixDefinition, 0)

	for _, def := range domain.AffixPool {
		if def.Type != affixType {
			continue
		}

		if ctx.Item.ItemLevel < def.MinLevel {
			continue
		}

		validAffixes = append(validAffixes, def)
	}

	if len(validAffixes) == 0 {
		return domain.AffixInstance{}, errors.New("No valid affixes")
	}

	chosenAffix := weightedRoll(validAffixes)

	value := ctx.RNG.Intn(chosenAffix.MaxValue-chosenAffix.MinValue+1) + chosenAffix.MinValue

	return domain.AffixInstance{
		DefID: chosenAffix.ID,
		Value: value,
	}, nil
}
