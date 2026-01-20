package domain

import (
	"SSVC-Server/internal/crafting/service"
	"errors"
)

type OrbOfTransmutation struct{}

func (o *OrbOfTransmutation) Apply(ctx *CraftingContext, affixType AffixType) error {
	item := ctx.Item

	if item.Rarity != Normal {
		return errors.New("Item must be normal")
	}

	item.Rarity = Magic

	switch affixType {
	case Prefix:
		service.RollAffix(ctx, Prefix)
	case Suffix:
		service.RollAffix(ctx, Suffix)
	}

	item.Prefixes = []AffixInstance{prefix}
	item.Suffixes = []AffixInstance{suffix}

	return nil
}
