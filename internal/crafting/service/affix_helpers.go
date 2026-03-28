package service

import "SSVC-Server/internal/crafting/domain"

func resolveSelectedAffixType(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) domain.AffixType {

	item := ctx.Item
	limits := domain.AffixLimitsByRarity[item.Rarity]

	switch affixType {

	case domain.Prefix, domain.OnlyPrefixes:
		return domain.Prefix

	case domain.Suffix, domain.OnlySuffixes:
		return domain.Suffix

	case domain.Both:
		prefixAvailable := len(item.Prefixes) < limits.MaxPrefixes
		suffixAvailable := len(item.Suffixes) < limits.MaxSuffixes

		switch {
		case prefixAvailable && suffixAvailable:
			if ctx.RNG.Intn(2) == 0 {
				return domain.Prefix
			}
			return domain.Suffix

		case prefixAvailable:
			return domain.Prefix

		default:
			return domain.Suffix
		}
	}

	return affixType
}
