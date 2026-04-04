package domain

import "testing"

func TestHasPrefixModifier(t *testing.T) {
	item := Item{
		Prefixes: []AffixDefinition{
			{ID: "lock_suffixes"},
		},
	}

	if !item.HasPrefixModifier("lock_suffixes") {
		t.Fatal("expected true")
	}
}

func TestHasSuffixModifier(t *testing.T) {
	item := Item{
		Suffixes: []AffixDefinition{
			{ID: "lock_prefixes"},
		},
	}

	if !item.HasSuffixModifier("lock_prefixes") {
		t.Fatal("expected true")
	}
}
