package service

import (
	"errors"
	"math"
	"testing"

	"SSVC-Server/internal/game/domain"
)

//Helpers

var rarityNames = map[domain.Rarity]string{
	domain.Normal: "Normal",
	domain.Magic:  "Magic",
	domain.Rare:   "Rare",
}

type mockRNG struct {
	value float64
}

func (m *mockRNG) Floatn(n float64) float64 {
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
		RNG:  &mockRNG{42},
	}
}

func TestWeightedRoll_EmptyPool(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on empty pool")
		}
	}()

	weightedRoll(&mockRNG{42}, []domain.BaseAffix{})
}

func TestWeightedRoll_SingleElement(t *testing.T) {
	def := domain.BaseAffix{Weight: 10}

	for i := 0; i < 100; i++ {
		result := weightedRoll(&mockRNG{42}, []domain.BaseAffix{def})
		if result.ID != def.ID {
			t.Fatalf("expected %v, got %v", def, result)
		}
	}
}

func TestWeightedRoll_Bias(t *testing.T) {
	pool := []domain.BaseAffix{
		{Weight: 1}, // rare
		{Weight: 9}, // common
	}

	counts := make(map[int]int)

	for i := 0; i < 10000; i++ {
		result := weightedRoll(&mockRNG{42}, pool)
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
				RNG: &mockRNG{float64(tc.mockRoll)},
			}

			step := RemoveIntegrity(int(math.Round(float64(tc.subtractRange))))
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

func TestAddIntegrity(t *testing.T) {
	tests := []struct {
		name           string
		startIntegrity int
		addAmount      int
		wantIntegrity  int
		wantErr        error
	}{
		{
			name:           "Add integrity below max",
			startIntegrity: 50,
			addAmount:      20,
			wantIntegrity:  70,
			wantErr:        nil,
		},
		{
			name:           "Add integrity reaching max",
			startIntegrity: 90,
			addAmount:      15,
			wantIntegrity:  100, // Max is 100
			wantErr:        nil,
		},
		{
			name:           "No addition if integrity is full",
			startIntegrity: 100,
			addAmount:      10,
			wantIntegrity:  100,
			wantErr:        errors.New("No item Integrity to add."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &domain.Item{Integrity: tt.startIntegrity}
			ctx := &domain.CraftingContext{Item: item}

			step := AddIntegrity(tt.addAmount)
			err := step(ctx)

			if (err != nil && tt.wantErr == nil) || (err == nil && tt.wantErr != nil) ||
				(err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) {
				t.Fatalf("unexpected error: got %v, want %v", err, tt.wantErr)
			}

			if ctx.Item.Integrity != tt.wantIntegrity {
				t.Errorf("unexpected integrity: got %d, want %d", ctx.Item.Integrity, tt.wantIntegrity)
			}
		})
	}
}

//Catalyst Tests

func TestImbuementCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Normal)

	c := &ImbuementCatalyst{}
	if err := c.Apply(ctx, domain.Either, &mockRNG{42}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Item.Rarity != domain.Magic {
		t.Fatalf("expected Magic rarity, got %v", ctx.Item.Rarity)
	}

	if len(ctx.Item.Prefixes)+len(ctx.Item.Suffixes) != 1 {
		t.Fatal("expected either 1 prefix or 1 suffix applied")
	}
}

func TestReconstructingCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Magic)

	c := &ReconstructingCatalyst{}
	if err := c.Apply(ctx, &mockRNG{42}); err != nil {
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
	ctx.Item.Prefixes = []domain.AffixDefinition{{ID: "test_prefix", Name: "test_prefix", Type: domain.Prefix, Tags: []string{"test_prefix"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{ID: "test_suffix", Name: "test_suffix", Type: domain.Suffix, Tags: []string{"test_suffix"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}}

	c := &ElevatingCatalyst{}
	if err := c.Apply(ctx, domain.Either, &mockRNG{42}); err != nil {
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

	domain.AffixPool = append(domain.AffixPool,
		domain.AffixDefinition{ID: "test_prefix", Name: "test_prefix", Type: domain.Prefix, Tags: []string{"test_prefix"}, MinValue: 1, MaxValue: 1, Weight: 100, MinLevel: 1},
		domain.AffixDefinition{ID: "test_suffix", Name: "test_suffix", Type: domain.Suffix, Tags: []string{"test_suffix"}, MinValue: 1, MaxValue: 1, Weight: 100, MinLevel: 1},
	)

	ctx.Item.Prefixes = []domain.AffixDefinition{{ID: "test_prefix", Name: "test_prefix", Type: domain.Prefix, Tags: []string{"test_prefix"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{ID: "test_suffix", Name: "test_suffix", Type: domain.Suffix, Tags: []string{"test_suffix"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}}

	oldAffixCount := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)

	c := &AscendantCatalyst{}

	if err := c.Apply(ctx, domain.Either, &mockRNG{42}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newAffixCount := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)

	if newAffixCount <= oldAffixCount {
		t.Fatalf("expected more affixes applied, got %d", newAffixCount)
	}

}

func TestDefiantCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	ctx.Item.Prefixes = []domain.AffixDefinition{}
	ctx.Item.Suffixes = []domain.AffixDefinition{}

	c := &DefiantCatalyst{}
	if err := c.Apply(ctx, &mockRNG{42}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Item.Prefixes) == 0 || len(ctx.Item.Suffixes) == 0 {
		t.Fatal("expected affixes to be applied")
	}

	if len(ctx.Item.Prefixes) < 1 || len(ctx.Item.Suffixes) < 1 {
		t.Fatal("expected between 1 and 3 affixes applied")
	}

	if len(ctx.Item.Prefixes) < 3 || len(ctx.Item.Suffixes) < 3 {
		t.Fatal("expected at least 3 affixes applied")
	}

	if len(ctx.Item.Prefixes) > 3 || len(ctx.Item.Suffixes) > 3 {
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

	ctx.Item.Prefixes = []domain.AffixDefinition{{ID: "test", Name: "Test", Type: domain.Prefix, Tags: []string{"test"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}, {ID: "test2", Name: "Test2", Type: domain.Suffix, Tags: []string{"test2"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}}
	ctx.Item.Suffixes = []domain.AffixDefinition{{ID: "lock_prefixes", Name: "lock_prefixes", Type: domain.Suffix, Tags: []string{"lock_prefixes"}, MinValue: 1, MaxValue: 1, DisplayedValues: []float64{1}, Weight: 100, MinLevel: 1}}

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

func TestCatharticCatalyst_Apply(t *testing.T) {

	ctx := newTestContext(domain.Rare)

	c := &CatharticCatalyst{}

	if err := c.Apply(ctx, &mockRNG{42}); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	if ctx.Item.Integrity > 60 || ctx.Item.Integrity < 50 {
		t.Errorf("Integrity after Apply = %d, want between 50 and 60", ctx.Item.Integrity)
	}

	totalAffixes := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)
	if totalAffixes < 1 || totalAffixes > 6 {
		t.Errorf("Total affixes after Apply = %d, want between 1 and 6", totalAffixes)
	}

	expectedRarity := InferRarityFromAffixes(ctx.Item)

	if ctx.Item.Rarity != expectedRarity {
		t.Errorf("Rarity after Apply = %v, want %v", rarityNames[ctx.Item.Rarity], rarityNames[expectedRarity])
	}
}

func TestCatharticCatalyst_Apply_WithLockedMods(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	ctx.Item.Prefixes = []domain.AffixDefinition{
		{
			ID:              "test_prefix",
			Name:            "Test Prefix",
			Type:            domain.Prefix,
			Tags:            []string{"test"},
			MinValue:        1,
			MaxValue:        1,
			DisplayedValues: []float64{1},
			Weight:          100,
			MinLevel:        1,
		},
		{
			ID:              "test_prefix_2",
			Name:            "Test Prefix 2",
			Type:            domain.Prefix,
			Tags:            []string{"test"},
			MinValue:        1,
			MaxValue:        1,
			DisplayedValues: []float64{1},
			Weight:          100,
			MinLevel:        1,
		},
	}
	ctx.Item.Suffixes = []domain.AffixDefinition{
		{
			ID:              "lock_prefixes",
			Name:            "Lock Prefixes",
			Type:            domain.Suffix,
			Tags:            []string{"lock_prefixes"},
			MinValue:        1,
			MaxValue:        1,
			DisplayedValues: []float64{1},
			Weight:          100,
			MinLevel:        1,
		},
	}

	c := &CatharticCatalyst{}

	if err := c.Apply(ctx, &mockRNG{42}); err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	if ctx.Item.Integrity > 60 || ctx.Item.Integrity < 50 {
		t.Errorf("Integrity after Apply = %d, want between 50 and 60", ctx.Item.Integrity)
	}

	totalAffixes := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)
	if totalAffixes < 2 || totalAffixes > 6 {
		t.Errorf("Total affixes after Apply = %d, want between 2 and 6", totalAffixes)
	}

	foundPrefix := false
	for _, p := range ctx.Item.Prefixes {
		if p.ID == "test_prefix" {
			foundPrefix = true
		}
	}

	if !foundPrefix {
		t.Errorf("Locked prefix was removed or changed: %+v", ctx.Item.Prefixes)
	}

	t.Logf("%v", ctx.Item)

	expectedRarity := InferRarityFromAffixes(ctx.Item)
	if ctx.Item.Rarity != expectedRarity {
		t.Errorf("Rarity after Apply = %v, want %v", rarityNames[ctx.Item.Rarity], rarityNames[expectedRarity])
	}
}
