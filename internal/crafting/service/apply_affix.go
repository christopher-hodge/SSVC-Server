package service

import "SSVC-Server/internal/crafting/domain"

func ApplyAffix(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) error {

	limits := domain.AffixLimitsByRarity[ctx.Item.Rarity]

	if err := domain.CanAddAffix(ctx.Item, affixType); err != nil {
		return err
	}

	//Only checks if the affixType is Both, then sets it to a vaild affix
	if affixType == domain.Both {
		switch {
		case len(ctx.Item.Prefixes) < limits.MaxPrefixes && len(ctx.Item.Suffixes) < limits.MaxSuffixes:
			affixType = domain.AffixType(ctx.RNG.Intn(2))

		case len(ctx.Item.Prefixes) < limits.MaxPrefixes:
			affixType = domain.Prefix

		default:
			affixType = domain.Suffix
		}
	}

	affix, err := RollAffix(ctx, affixType)
	if err != nil {
		return err
	}

	item := ctx.Item

	if affixType == domain.Prefix {
		item.Prefixes = append(item.Prefixes, affix)
	} else {
		item.Suffixes = append(item.Suffixes, affix)
	}

	return nil
}
