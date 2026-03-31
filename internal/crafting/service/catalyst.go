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

func ClearAffixes() CraftStep {
	return func(ctx *domain.CraftingContext) error {

		if !ctx.Item.HasPrefixModifier("lock_prefixes") {
			ctx.Item.Suffixes = nil
		}

		if !ctx.Item.HasSuffixModifier("lock_suffixes") {
			ctx.Item.Prefixes = nil
		}

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
		AddAffixes(1, domain.Both),
	})
}

func (c *ReconstructionCatalyst) Apply(ctx *domain.CraftingContext) error {
	count := ctx.RNG.Intn(2) + 1

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		ClearAffixes(),
		AddAffixes(count, domain.All),
	})
}

func (c *ElevatingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		SetRarity(domain.Rare),
		AddAffixes(1, domain.Both),
	})
}

func (c *DefiantCatalyst) Apply(ctx *domain.CraftingContext) error {
	count := ctx.RNG.Intn(6) + 1

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		ClearAffixes(),
		AddAffixes(count, domain.All),
	})
}

func (c *AscendantCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		AddAffixes(1, domain.Both),
	})
}

func (c *LustratingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {

	newRarity := InferRarityFromAffixes(ctx.Item)

	return ExecutePipeline(ctx, []CraftStep{
		ClearAffixes(),
		SetRarity(newRarity),
	})
}
