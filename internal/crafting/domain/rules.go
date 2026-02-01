package domain

import "errors"

func CanAddAffix(item *Item, affixType AffixType) error {

	maxPrefixes := 0
	maxSuffixes := 0

	switch item.Rarity {
	case Normal:
		maxPrefixes = 0
		maxSuffixes = 0
	case Magic:
		maxPrefixes = 1
		maxSuffixes = 1
	case Rare:
		maxPrefixes = 3
		maxSuffixes = 3
	}

	if affixType == Prefix && len(item.Prefixes) >= maxPrefixes {
		return errors.New("max prefixes reached")
	}

	if affixType == Suffix && len(item.Suffixes) >= maxSuffixes {
		return errors.New("max suffixes reached")
	}

	return nil
}
