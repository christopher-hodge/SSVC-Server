package service

import (
	"SSVC-Server/internal/game/domain"
	"SSVC-Server/internal/random"
	"errors"
	"math"
)

func ApplyAffixLogic(ctx *domain.CraftingContext, affixType domain.AffixType, rng random.RNGerFloat) error {

	limits := domain.AffixLimitsByRarity[ctx.Item.Rarity]
	prefixAvailable := len(ctx.Item.Prefixes) < limits.MaxPrefixes
	suffixAvailable := len(ctx.Item.Suffixes) < limits.MaxSuffixes

	if affixType == domain.Either {

		switch {
		case prefixAvailable && suffixAvailable:
			affixType = domain.Prefix
			if ctx.RNG.Floatn(math.Round(2)) == 1 {
				affixType = domain.Suffix
			}
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

	affix, err := RollAffix(ctx, rng, affixType)
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
	rng random.RNGerFloat,
	affixType domain.AffixType,
) (domain.AffixDefinition, error) {

	validAffixes := make([]domain.BaseAffix, 0)

	for _, base := range domain.BaseAffixes {
		if base.Type != affixType {
			continue
		}

		if ctx.Item.ItemLevel < base.MinLevel {
			continue
		}

		validAffixes = append(validAffixes, base)
	}

	if len(validAffixes) == 0 {
		return domain.AffixDefinition{}, errors.New("no valid affixes")
	}

	// 1️⃣ pick base affix
	chosenBase := weightedRoll(rng, validAffixes)

	// 2️⃣ roll tier (you can improve this later)
	tier := rollTier(ctx, chosenBase, rng)

	// 3️⃣ generate full affix (handles scaling + values)
	affix := domain.GenerateAffix(chosenBase, tier, rng)

	return affix, nil
}

func rollTier(
	ctx *domain.CraftingContext,
	base domain.BaseAffix,
	rng random.RNGerFloat,
) int {

	maxTier := 10

	tier := rng.Floatn(float64(maxTier)) + 1

	return int(tier)
}

func weightedRoll(
	rng random.RNGerFloat,
	pool []domain.BaseAffix,
) domain.BaseAffix {

	totalWeight := 0
	for _, def := range pool {
		totalWeight += def.Weight
	}

	roll := rng.Floatn(float64(totalWeight))

	running := 0
	for _, def := range pool {
		running += def.Weight
		if roll < float64(running) {
			return def
		}
	}

	return pool[len(pool)-1]
}
