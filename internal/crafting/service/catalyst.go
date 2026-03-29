package service

import (
	"errors"

	"SSVC-Server/internal/crafting/domain"
)

type ImbuementCatalyst struct{}
type ReconstructionCatalyst struct{}
type ElevatingCatalyst struct{}
type DefiantCatalyst struct{}
type AscendantCatalyst struct{}
type LustratingCatalyst struct{}

type CraftStep func(ctx *domain.CraftingContext) error

func ExecutePipeline(ctx *domain.CraftingContext, steps []CraftStep) error {
	for _, step := range steps {
		if err := step(ctx); err != nil {
			return err
		}
	}
	return nil
}

func RequireRarity(r domain.Rarity) CraftStep {
	return func(ctx *domain.CraftingContext) error {
		if ctx.Item.Rarity != r {
			return errors.New("Invalid item rarity.")
		}
		return nil
	}
}

func ClearAffixes(i domain.Item) CraftStep {
	if !i.HasPrefixModifier("lock_prefixes") {
		return func(ctx *domain.CraftingContext) error {
			i.Suffixes = []domain.AffixInstance{}
			return nil
		}
	}
	if !i.HasSuffixModifier("lock_suffixes") {
		return func(ctx *domain.CraftingContext) error {
			ctx.Item.Prefixes = []domain.AffixInstance{}
			return nil
		}
	}

	return func(ctx *domain.CraftingContext) error {
		ctx.Item.Prefixes = []domain.AffixInstance{}
		ctx.Item.Suffixes = []domain.AffixInstance{}
		return nil
	}
}

func SetRarity(r domain.Rarity) CraftStep {
	return func(ctx *domain.CraftingContext) error {
		ctx.Item.Rarity = r
		return nil
	}
}

func InferRarityFromAffixes(item *domain.Item) domain.Rarity {
	prefixCount := len(item.Prefixes)
	suffixCount := len(item.Suffixes)

	if prefixCount == 0 && suffixCount == 0 {
		return domain.Normal
	}

	magic := domain.AffixLimitsByRarity[domain.Magic]
	if prefixCount <= magic.MaxPrefixes && suffixCount <= magic.MaxSuffixes {
		return domain.Magic
	}

	return domain.Rare
}

func AddAffix(affixType domain.AffixType) CraftStep {
	return func(ctx *domain.CraftingContext) error {
		return ApplyAffixLogic(ctx, affixType)
	}
}

func AddAffixes(count int, affixType domain.AffixType) CraftStep {
	return func(ctx *domain.CraftingContext) error {
		for i := 0; i < count; i++ {
			if err := ApplyAffixLogic(ctx, affixType); err != nil {
				return err
			}
		}
		return nil
	}
}

//Crafting item functions

func (c *ImbuementCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Normal),
		SetRarity(domain.Magic),
		AddAffix(domain.Both),
	})
}

func (c *ReconstructionCatalyst) Apply(ctx *domain.CraftingContext) error {
	count := ctx.RNG.Intn(2) + 1 // 1 or 2 mods

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		ClearAffixes(*ctx.Item),
		AddAffixes(count, domain.All),
	})
}

func (c *ElevatingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		SetRarity(domain.Rare),
		AddAffix(domain.Both),
	})
}

func (c *DefiantCatalyst) Apply(ctx *domain.CraftingContext) error {
	count := ctx.RNG.Intn(6) + 1

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		ClearAffixes(*ctx.Item),
		AddAffixes(count, domain.All),
	})
}

func (c *AscendantCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		AddAffix(domain.Both),
	})
}

func (c *LustratingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		ClearAffixes(*ctx.Item),
		SetRarity(InferRarityFromAffixes(ctx.Item)),
	})
}
