package domain

import (
	"fmt"
	"testing"
)

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

func TestGenerateAffixPool(t *testing.T) {
	InitAffixPool() // populate AffixPool

	if len(AffixPool) == 0 {
		t.Fatal("AffixPool is empty, affixes not generated")
	}

	for _, a := range AffixPool {
		fmt.Printf("%s: %d-%d, Weight %d\n", a.Name, a.MinValue, a.MaxValue, a.Weight)
	}
}
