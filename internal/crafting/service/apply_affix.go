package service

import (
	"SSVC-Server/internal/crafting/domain"
	"errors"
)

func ApplyAffixLogic(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	limits := domain.AffixLimitsByRarity[ctx.Item.Rarity]

	if affixType == domain.Prefix && ctx.Item.HasPrefixModifier("lock_suffixes") {
		return errors.New("cannot roll suffixes")
	}

	if affixType == domain.Suffix && ctx.Item.HasSuffixModifier("lock_prefixes") {
		return errors.New("cannot roll prefixes")
	}

	if affixType == domain.Both || affixType == domain.All {
		prefixAvailable := len(ctx.Item.Prefixes) < limits.MaxPrefixes
		suffixAvailable := len(ctx.Item.Suffixes) < limits.MaxSuffixes

		switch {
		case prefixAvailable && suffixAvailable:
			affixType = domain.AffixType(ctx.RNG.Intn(2)) // pick randomly
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

	affix, err := RollAffix(ctx, affixType)
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
