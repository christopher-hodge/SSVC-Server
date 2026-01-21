package domain

import (
	"errors"
	"math/rand"
)

type ImbuementCatalyst struct{}
type ReconstructionCatalyst struct{}

type AffixRoller interface {
	RollAffix(ctx *CraftingContext, affixType AffixType) (AffixInstance, error)
}

func rollAffixHelper(ctx *CraftingContext, item *Item, affixType AffixType) error {
	affix, err := ctx.AffixRoller.RollAffix(ctx, affixType)
	if err != nil {
		return err
	}

	if affixType == Prefix {
		item.Prefixes = append(item.Prefixes, affix)
	} else {
		item.Suffixes = append(item.Suffixes, affix)
	}

	return nil
}

func (o *ImbuementCatalyst) Apply(ctx *CraftingContext, affixType AffixType) error {
	item := ctx.Item

	if item.Rarity != Normal {
		return errors.New("Item must be Normal.")
	}

	item.Rarity = Magic

	return rollAffixHelper(ctx, item, affixType)
}

func (o *ReconstructionCatalyst) Apply(ctx *CraftingContext) error {
	item := ctx.Item

	if item.Rarity != Magic {
		return errors.New("Item must be Magic.")
	}

	outcomes := rand.Intn(3) // 0 - prefix, 1 - suffix, 2 - both

	switch outcomes {
	case 0:
		return rollAffixHelper(ctx, item, Prefix)
	case 1:
		return rollAffixHelper(ctx, item, Suffix)
	case 2:
		if err := rollAffixHelper(ctx, item, Prefix); err != nil {
			return err
		}
		return rollAffixHelper(ctx, item, Suffix)
	}

	return nil
}
