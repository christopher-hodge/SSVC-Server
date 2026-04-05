package service

import (
	"errors"
	"math"

	"SSVC-Server/internal/game/domain"
	"SSVC-Server/internal/random"
)

// Base Catalyst Structs
type ImbuementCatalyst struct{}
type ReconstructingCatalyst struct{}
type ElevatingCatalyst struct{}
type DefiantCatalyst struct{}
type AscendantCatalyst struct{}
type LustratingCatalyst struct{}

// Special Catalyst Structs
type CatharticCatalyst struct{}

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

		if ctx.Item.Integrity <= 0 {
			return errors.New("No integrity to remove.")
		}

		amount := ctx.RNG.Floatn(math.Round(float64(subtractRange))) + 1

		ctx.Item.Integrity -= int(amount)

		if ctx.Item.Integrity < 0 {
			ctx.Item.Integrity = 0
		}

		return nil
	}
}

func AddIntegrity(sumAmmount int) CraftStep {
	return func(ctx *domain.CraftingContext) error {

		if ctx.Item.Integrity == 100 {
			return errors.New("No item Integrity to add.")
		}

		if ctx.Item.Integrity+sumAmmount >= 100 {
			ctx.Item.Integrity = 100
			return nil
		} else {
			ctx.Item.Integrity = ctx.Item.Integrity + sumAmmount
			return nil
		}
	}
}

func AddAffixes(count int, affixType domain.AffixType, rng random.RNGerFloat) CraftStep {
	return func(ctx *domain.CraftingContext) error {
		for i := 0; i < count; i++ {
			if err := ApplyAffixLogic(ctx, affixType, rng); err != nil {
				break
			}
		}
		return nil
	}
}

//Base Crafting item functions

func (c *ImbuementCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType, rng random.RNGerFloat) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Normal),
		SetRarity(domain.Magic),
		RemoveIntegrity(5),
		AddAffixes(1, domain.Either, rng),
	})
}

func (c *ReconstructingCatalyst) Apply(ctx *domain.CraftingContext, rng random.RNGerFloat) error {
	count := ctx.RNG.Floatn(math.Round(2)) + 1

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		ClearAffixes(),
		RemoveIntegrity(5),
		AddAffixes(int(count), domain.Either, rng),
	})
}

func (c *ElevatingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType, rng random.RNGerFloat) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Magic),
		SetRarity(domain.Rare),
		RemoveIntegrity(10),
		AddAffixes(1, domain.Either, rng),
	})
}

func (c *DefiantCatalyst) Apply(ctx *domain.CraftingContext, rng random.RNGerFloat) error {

	count := ctx.RNG.Floatn(math.Round(4)) + 3 // Rares require at least 3 mods

	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		ClearAffixes(),
		RemoveIntegrity(15),
		AddAffixes(int(count), domain.Either, rng),
	})
}

func (c *AscendantCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType, rng random.RNGerFloat) error {
	return ExecutePipeline(ctx, []CraftStep{
		RequireRarity(domain.Rare),
		RemoveIntegrity(15),
		AddAffixes(1, domain.Either, rng),
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

func (c *CatharticCatalyst) Apply(ctx *domain.CraftingContext, rng random.RNGerFloat) error {

	count := ctx.RNG.Floatn(math.Round(7))

	return ExecutePipeline(ctx, []CraftStep{
		ClearAffixes(),
		RemoveIntegrity(10),
		AddAffixes(int(count), domain.Either, rng),
		func(ctx *domain.CraftingContext) error {
			ctx.Item.Rarity = InferRarityFromAffixes(ctx.Item)
			return nil
		},
	})
}
