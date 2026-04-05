package service

import (
	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/random"
	"errors"
)

func ApplyAffixLogic(ctx *domain.CraftingContext, affixType domain.AffixType) error {

	limits := domain.AffixLimitsByRarity[ctx.Item.Rarity]
	prefixAvailable := len(ctx.Item.Prefixes) < limits.MaxPrefixes
	suffixAvailable := len(ctx.Item.Suffixes) < limits.MaxSuffixes

	if affixType == domain.Prefix && ctx.Item.HasPrefixModifier("lock_suffixes") {
		return errors.New("cannot roll suffixes")
	}

	if affixType == domain.Suffix && ctx.Item.HasSuffixModifier("lock_prefixes") {
		return errors.New("cannot roll prefixes")
	}

	if affixType == domain.Either {

		switch {
		case prefixAvailable && suffixAvailable:
			affixType = domain.AffixType(ctx.RNG.Intn(1)) // pick randomly
		case prefixAvailable:
			affixType = domain.Prefix
		case suffixAvailable:
			affixType = domain.Suffix
		default:
			return errors.New("no affix slots available")
		}
	}

	if affixType == domain.Prefix && len(ctx.Item.Prefixes) >= limits.MaxPrefixes {
		return errors.New("max prefixes reached")
	}
	if affixType == domain.Suffix && len(ctx.Item.Suffixes) >= limits.MaxSuffixes {
		return errors.New("max suffixes reached")
	}

	affix, err := RollAffix(ctx, ctx.RNG, affixType)
	if err != nil {
		return err
	}

	if affixType == domain.Prefix {
		ctx.Item.Prefixes = append(ctx.Item.Prefixes, affix)
	} else {
		ctx.Item.Suffixes = append(ctx.Item.Suffixes, affix)
	}

	return nil
}

func RollAffix(
	ctx *domain.CraftingContext,
	rng random.RNGer,
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

	chosenAffix := weightedRoll(rng, validAffixes)

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
