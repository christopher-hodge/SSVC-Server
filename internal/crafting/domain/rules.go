package domain

import "errors"

func CanAddAffix(item *Item, affixType AffixType) error {

	if affixType == Prefix && len(item.Prefixes) >= 3 {
		return errors.New("max prefixes reached")
	}

	if affixType == Suffix && len(item.Suffixes) >= 3 {
		return errors.New("max suffixes reached")
	}

	return nil
}
