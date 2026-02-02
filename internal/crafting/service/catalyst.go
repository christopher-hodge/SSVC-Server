package service

import (
	"errors"

	"SSVC-Server/internal/crafting/domain"
)

type ImbuementCatalyst struct{}
type ReconstructionCatalyst struct{}
type ElevatingCatalyst struct{}

func (c *ImbuementCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {

	item := ctx.Item

	if item.Rarity != domain.Normal {
		return errors.New("Item must be Normal.")
	}

	item.Rarity = domain.Magic

	// Determines if this is a prefix (0) or suffix (1)
	newAffixType := ctx.RNG.Intn(2)

	return ApplyAffix(ctx, domain.AffixType(newAffixType))
}

func (c *ReconstructionCatalyst) Apply(ctx *domain.CraftingContext) error {

	item := ctx.Item

	if item.Rarity != domain.Normal && item.Rarity != domain.Magic {
		return errors.New("Item must be Normal or Magic.")
	}

	if item.Rarity == domain.Normal {
		item.Rarity = domain.Magic
	}

	item.Prefixes = []domain.AffixInstance{}
	item.Suffixes = []domain.AffixInstance{}

	//Determines if this is a prefix (0), suffix (1), or both(2)
	newAffixType := ctx.RNG.Intn(3)

	return ApplyAffix(ctx, domain.AffixType(newAffixType))
}

func (c *ElevatingCatalyst) Apply(ctx *domain.CraftingContext, affixType domain.AffixType) error {

	item := ctx.Item

	if item.Rarity != domain.Magic {
		return errors.New("Item must be Magic.")
	}

	item.Rarity = domain.Rare

	// Determines if this is a prefix (0) or suffix (1)
	newAffixType := ctx.RNG.Intn(2)

	return ApplyAffix(ctx, domain.AffixType(newAffixType))
}
