package service

import (
	"errors"
	"testing"

	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/random"
)

//Helpers

var rarityNames = map[domain.Rarity]string{
	domain.Normal: "Normal",
	domain.Magic:  "Magic",
	domain.Rare:   "Rare",
}

type mockRNG struct {
	value int
}

func (m *mockRNG) Intn(n int) int {
	return m.value
}

func newTestItem(rarity domain.Rarity) *domain.Item {
	return &domain.Item{
		Rarity:    rarity,
		ItemLevel: 71,
		Integrity: 100,
		Prefixes:  []domain.AffixDefinition{},
		Suffixes:  []domain.AffixDefinition{},
	}
}

func newTestContext(rarity domain.Rarity) *domain.CraftingContext {
	return &domain.CraftingContext{
		Item: newTestItem(rarity),
		RNG:  &mockRNG{value: 42},
	}
}

func TestWeightedRoll_EmptyPool(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty pool")
		}
	}()

	weightedRoll(random.New(42), []domain.AffixDefinition{})
}

func TestWeightedRoll_SingleElement(t *testing.T) {
	def := domain.AffixDefinition{Weight: 10}

	for i := 0; i < 100; i++ {
		result := weightedRoll(random.New(42), []domain.AffixDefinition{def})
		if result.ID != def.ID {
			t.Fatalf("expected %v, got %v", def, result)
		}
	}
}

func TestWeightedRoll_Bias(t *testing.T) {
	pool := []domain.AffixDefinition{
		{Weight: 1}, // rare
		{Weight: 9}, // common
	}

	counts := make(map[int]int)

	for i := 0; i < 10000; i++ {
		result := weightedRoll(random.New(42), pool)
		counts[result.Weight]++
	}

	if counts[9] <= counts[1] {
		t.Fatalf("expected higher weight to be rolled more often: %+v", counts)
	}
}

//Pipeline Tests

func TestExecutePipeline_Success(t *testing.T) {
	ctx := newTestContext(domain.Normal)

	steps := []CraftStep{
		func(ctx *domain.CraftingContext) error { return nil },
		func(ctx *domain.CraftingContext) error { return nil },
	}

	if err := ExecutePipeline(ctx, steps); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestExecutePipeline_FailureStopsExecution(t *testing.T) {
	ctx := newTestContext(domain.Normal)
	called := false

	steps := []CraftStep{
		func(ctx *domain.CraftingContext) error { return errors.New("invalid rarity") },
		func(ctx *domain.CraftingContext) error {
			called = true
			return nil
		},
	}

	if err := ExecutePipeline(ctx, steps); err == nil {
		t.Fatal("expected error, got nil")
	}
	if called {
		t.Fatal("expected pipeline to stop on error")
	}
}

func TestRemoveIntegrity(t *testing.T) {
	type testCase struct {
		name           string
		startIntegrity int
		subtractRange  int
		mockRoll       int
		expected       int
		expectError    bool
	}

	tests := []testCase{
		{
			name:           "normal reduction",
			startIntegrity: 10,
			subtractRange:  5,
			mockRoll:       2, // +1 = 3 removed
			expected:       7,
			expectError:    false,
		},
		{
			name:           "clamps to zero",
			startIntegrity: 1,
			subtractRange:  5,
			mockRoll:       4, // +1 = 5 removed
			expected:       0,
			expectError:    false,
		},
		{
			name:           "no integrity to remove",
			startIntegrity: 0,
			subtractRange:  5,
			mockRoll:       0,
			expected:       0,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &domain.CraftingContext{
				Item: &domain.Item{
					Integrity: tc.startIntegrity,
				},
				RNG: &mockRNG{value: tc.mockRoll},
			}

			step := RemoveIntegrity(tc.subtractRange)
			err := step(ctx)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ctx.Item.Integrity != tc.expected {
				t.Fatalf(
					"expected integrity %d, got %d",
					tc.expected,
					ctx.Item.Integrity,
				)
			}
		})
	}
}

//Catalyst Tests

func TestImbuementCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Normal)

	c := &ImbuementCatalyst{}
	if err := c.Apply(ctx, domain.Either); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Item.Rarity != domain.Magic {
		t.Fatalf("expected Magic rarity, got %v", ctx.Item.Rarity)
	}

	if len(ctx.Item.Prefixes)+len(ctx.Item.Suffixes) != 1 {
		t.Fatal("expected either 1 prefix or 1 suffix applied")
	}
}

func TestReconstructionCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Magic)

	c := &ReconstructingCatalyst{}
	if err := c.Apply(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Item.Rarity != domain.Magic {
		t.Fatalf("expected Magic rarity, got %v", ctx.Item.Rarity)
	}

	if len(ctx.Item.Prefixes) == 0 || len(ctx.Item.Suffixes) == 0 {
		t.Fatal("expected affixes to be applied")
	}

	if len(ctx.Item.Prefixes) > 1 || len(ctx.Item.Suffixes) > 1 {
		t.Fatal("too many affixes applied, expected 1 prefix and 1 suffix")
	}
}

func TestElevatingCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Magic)
	ctx.Item.Prefixes = []domain.AffixDefinition{{ID: "test_prefix", Name: "test_prefix", Type: domain.Prefix, Tags: []string{"test_prefix"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{ID: "test_suffix", Name: "test_suffix", Type: domain.Suffix, Tags: []string{"test_suffix"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}}

	c := &ElevatingCatalyst{}
	if err := c.Apply(ctx, domain.Either); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Item.Rarity != domain.Rare {
		t.Fatalf("expected Rare rarity, got %v", ctx.Item.Rarity)
	}

	totalAffixCount := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)

	if totalAffixCount < 3 {
		t.Fatalf("expected at least 3 affixes applied got %d", totalAffixCount)
	}
}

func TestAscendantCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)
	ctx.Item.Prefixes = []domain.AffixDefinition{{ID: "test_prefix", Name: "test_prefix", Type: domain.Prefix, Tags: []string{"test_prefix"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{ID: "test_suffix", Name: "test_suffix", Type: domain.Suffix, Tags: []string{"test_suffix"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}}

	oldAffixCount := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)

	c := &AscendantCatalyst{}

	if err := c.Apply(ctx, domain.Either); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newAffixCount := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)

	if newAffixCount <= oldAffixCount {
		t.Fatalf("expected more affixes applied, got %d", newAffixCount)
	}

}

func TestDefiantCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	c := &DefiantCatalyst{}
	if err := c.Apply(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Item.Prefixes) == 0 || len(ctx.Item.Suffixes) == 0 {
		t.Fatal("expected affixes to be applied")
	}

	if len(ctx.Item.Prefixes) < 1 || len(ctx.Item.Suffixes) < 1 {
		t.Fatal("expected between 1 and 6 affixes applied")
	}

	if len(ctx.Item.Prefixes) < 3 || len(ctx.Item.Suffixes) < 3 {
		t.Fatal("expected at least 3 affixes applied")
	}

	if len(ctx.Item.Prefixes) > 6 || len(ctx.Item.Suffixes) > 6 {
		t.Fatal("too many affixes applied, expected between 1 and 6")
	}
}

func TestLustratingCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	ctx.Item.Prefixes = []domain.AffixDefinition{{}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{}}

	c := &LustratingCatalyst{}
	if err := c.Apply(ctx, domain.Either); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Item.Prefixes) == 0 && len(ctx.Item.Suffixes) == 0 {
	} else {
		t.Fatal("expected all affixes cleared")
	}

	if ctx.Item.Rarity != domain.Normal {
		t.Fatalf("expected Normal rarity after clearing, got %s", rarityNames[ctx.Item.Rarity])
	}
}

func TestLustratingCatalyst_WithLockedMods(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	ctx.Item.Prefixes = []domain.AffixDefinition{{ID: "test", Name: "Test", Type: domain.Prefix, Tags: []string{"test"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}, {ID: "test2", Name: "Test2", Type: domain.Suffix, Tags: []string{"test2"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{ID: "lock_prefixes", Name: "lock_prefixes", Type: domain.Suffix, Tags: []string{"lock_prefixes"}, MinValue: 1, MaxValue: 1, DisplayedValue: 1, Weight: 100, MinLevel: 1}}

	c := &LustratingCatalyst{}
	if err := c.Apply(ctx, domain.Either); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Item.Prefixes) != 2 {
		t.Fatalf("expected prefixes to remain, got %d", len(ctx.Item.Prefixes))
	}

	if len(ctx.Item.Suffixes) != 0 {
		t.Fatalf("expected suffixes cleared, got %d", len(ctx.Item.Suffixes))
	}

	if ctx.Item.Rarity != domain.Rare {
		t.Fatalf("expected Rare rarity after clearing, got %s", rarityNames[ctx.Item.Rarity])
	}
}
