package service

import (
	"errors"

	"SSVC-Server/internal/crafting/domain"
)

// Base Catalyst Structs
type ImbuementCatalyst struct{}
type ReconstructingCatalyst struct{}
type ElevatingCatalyst struct{}
type DefiantCatalyst struct{}
type AscendantCatalyst struct{}
type LustratingCatalyst struct{}

// Special Catalyst Structs
type CatharsisCatalyst struct{}

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

		if !ctx.Item.HasSuffixModifier("lock_prefixes") {
			ctx.Item.Prefixes = nil
		}

		if !ctx.Item.HasPrefixModifier("lock_suffixes") {
			ctx.Item.Suffixes = nil
		}

		return nil
	}
}

func SetRarity(rarity domain.Rarity) CraftStep {
	return func(ctx *domain.CraftingContext) error {
		ctx.Item.Rarity = rarity
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

func InferAndSetRarity() CraftStep {
	return func(ctx *domain.CraftingContext) error {
		ctx.Item.Rarity = InferRarityFromAffixes(ctx.Item)
		return nil
	}
}

func RemoveIntegrity(subtractRange int) CraftStep {
	return func(ctx *domain.CraftingContext) error {

		if ctx.Item.Integrity != 0 {
			ctx.Item.Integrity = ctx.Item.Integrity - ctx.RNG.Intn(subtractRange) + 1
			return nil
		} else {
			return errors.New("No item Integrity to remove.")
		}
	}
}

func AddIntegrity(sumAmmount int) CraftStep {
	return func(ctx *domain.CraftingContext) error {

		if ctx.Item.Integrity == 100 {
			return errors.New("No item Integrity to add.")
		}

		if ctx.Item.Integrity+sumAmmount >= 100 {
			ctx.Item.Integrity = ctx.Item.Integrity + sumAmmount
			return nil
		} else {
			ctx.Item.Integrity = ctx.Item.Integrity + sumAmmount
			return nil
		}
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

//Base Crafting item functions

func (c *ImbuementCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Normal),
		SetRarity(domain.Magic),
		RemoveIntegrity(5),
		AddAffixes(1, domain.Either),
	})
}

func (c *ReconstructingCatalyst) Apply(ctx *domain.CraftingContext) error {
	count := ctx.RNG.Intn(2) + 1

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		ClearAffixes(),
		RemoveIntegrity(5),
		AddAffixes(count, domain.Either),
	})
}

func (c *ElevatingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		SetRarity(domain.Rare),
		RemoveIntegrity(10),
		AddAffixes(1, domain.Either),
	})
}

func (c *DefiantCatalyst) Apply(ctx *domain.CraftingContext) error {

	count := ctx.RNG.Intn(6) + 1

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		ClearAffixes(),
		RemoveIntegrity(15),
		AddAffixes(count, domain.Either),
	})
}

func (c *AscendantCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		RemoveIntegrity(15),
		AddAffixes(1, domain.Either),
	})
}

func (c *LustratingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {

	return ExecutePipeline(ctx, []CraftStep{
		ClearAffixes(),
		RemoveIntegrity(5),
		InferAndSetRarity(),
	})
}

//Special Catalyst Crafting functions

func (c *CatharsisCatalyst) Apply(ctx *domain.CraftingContext) error {

	count := ctx.RNG.Intn(6) + 1

	return ExecutePipeline(ctx, []CraftStep{
		ClearAffixes(),
		RemoveIntegrity(10),
		AddAffixes(count, domain.Either),
		SetRarity(InferRarityFromAffixes(ctx.Item)),
	})

}
