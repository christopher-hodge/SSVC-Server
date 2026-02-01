package service

import "SSVC-Server/internal/crafting/domain"

func ApplyAffix(
	ctx *domain.CraftingContext,
	affixType domain.AffixType,
) error {

	if err := domain.CanAddAffix(ctx.Item, affixType); err != nil {
		return err
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
