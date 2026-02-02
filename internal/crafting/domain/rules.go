package domain

import "errors"

func CanAddAffix(item *Item, affixType AffixType) error {

	if item.Rarity == Normal {
		return errors.New("Normal items cannot have affixes.")
	}

	limits := AffixLimitsByRarity[item.Rarity]

	switch affixType {
	case Prefix:
		if len(item.Prefixes) >= limits.MaxPrefixes {
			return errors.New("Maximum prefixes reached.")
		}

	case Suffix:
		if len(item.Suffixes) >= limits.MaxSuffixes {
			return errors.New("Maximum suffixes reached.")
		}

	case Both:
		prefixFull := len(item.Prefixes) >= limits.MaxPrefixes
		suffixFull := len(item.Suffixes) >= limits.MaxSuffixes

		if prefixFull && suffixFull {
			return errors.New("Maximum affixes reached.")
		}

	}

	return nil
}

func CanRemoveAffix(item *Item, affixType AffixType) error {

	switch affixType {
	case Prefix:
		if len(item.Prefixes) == 0 {
			return errors.New("No prefixes are able to be removed.")
		}
	case Suffix:
		if len(item.Suffixes) == 0 {
			return errors.New("No suffixes are able to be removed.")
		}
	}

	return nil
}
