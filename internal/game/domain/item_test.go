package domain

import (
	"fmt"
	"testing"
)

type mockRNG struct {
	value int
}

func (m *mockRNG) Intn(n int) int {
	return m.value
}

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
	InitAffixPool(&mockRNG{42})

	if len(AffixPool) == 0 {
		t.Fatal("AffixPool is empty, affixes not generated")
	}

	for _, a := range AffixPool {
		fmt.Printf("%s: %d-%d, Weight %d\n", a.Name, a.MinValue, a.MaxValue, a.Weight)
	}
}
