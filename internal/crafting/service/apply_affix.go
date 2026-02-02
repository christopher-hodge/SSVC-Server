package service

import "SSVC-Server/internal/crafting/domain"

func ApplyAffix(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) error {

	maxPrefixes := 0
	maxSuffixes := 0

	switch ctx.Item.Rarity {
	case domain.Magic:
		maxPrefixes = 1
		maxSuffixes = 1
	case domain.Rare:
		maxPrefixes = 3
		maxSuffixes = 3
	}

	if err := domain.CanAddAffix(ctx.Item, affixType); err != nil {
		return err
	}

	//Only checks if the affixType is Both, then sets it to a vaild affix
	if affixType == domain.Both {
		switch {
		case len(ctx.Item.Prefixes) < maxPrefixes && len(ctx.Item.Suffixes) < maxSuffixes:
			affixType = domain.AffixType(ctx.RNG.Intn(2))

		case len(ctx.Item.Prefixes) < maxPrefixes:
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
