package service

import (
	"errors"

	"SSVC-Server/internal/crafting/domain"
)

func RollAffix(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) (domain.AffixDefinition, error) {

	validAffixes := make([]domain.AffixDefinition, 0)

	for _, def := range domain.AffixPool {
		if def.Type != affixType {
			continue
		}

		if ctx.Item.ItemLevel <= def.MinLevel {
			continue
		}

		validAffixes = append(validAffixes, def)
	}

	if len(validAffixes) == 0 {
		return domain.AffixDefinition{}, errors.New("No valid affixes")
	}

	chosenAffix := weightedRoll(validAffixes)

	value := ctx.RNG.Intn(chosenAffix.MaxValue-chosenAffix.MinValue+1) + chosenAffix.MinValue

	return domain.AffixDefinition{
		ID:             chosenAffix.ID,
		Name:           chosenAffix.Name,
		Type:           chosenAffix.Type,
		Tags:           chosenAffix.Tags,
		MinValue:       chosenAffix.MinValue,
		MaxValue:       chosenAffix.MaxValue,
		DisplayedValue: value,
		Weight:         chosenAffix.Weight,
		MinLevel:       chosenAffix.MinLevel,
	}, nil
}
