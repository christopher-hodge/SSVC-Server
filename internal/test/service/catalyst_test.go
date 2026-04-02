package service

import (
	"errors"
	"testing"

	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/crafting/service"
	"SSVC-Server/internal/random"
)

//Helpers

var rarityNames = map[domain.Rarity]string{
	domain.Normal: "Normal",
	domain.Magic:  "Magic",
	domain.Rare:   "Rare",
}

func newTestItem(rarity domain.Rarity) *domain.Item {
	return &domain.Item{
		Rarity:    rarity,
		ItemLevel: 71,
		Prefixes:  []domain.AffixDefinition{},
		Suffixes:  []domain.AffixDefinition{},
	}
}

func newTestContext(rarity domain.Rarity) *domain.CraftingContext {
	return &domain.CraftingContext{
		Item: newTestItem(rarity),
		RNG:  *random.New(42),
	}
}

//Pipeline Tests

func TestExecutePipeline_Success(t *testing.T) {
	ctx := newTestContext(domain.Normal)

	steps := []service.CraftStep{
		func(ctx *domain.CraftingContext) error { return nil },
		func(ctx *domain.CraftingContext) error { return nil },
	}

	if err := service.ExecutePipeline(ctx, steps); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestExecutePipeline_FailureStopsExecution(t *testing.T) {
	ctx := newTestContext(domain.Normal)
	called := false

	steps := []service.CraftStep{
		func(ctx *domain.CraftingContext) error { return errors.New("invalid rarity") },
		func(ctx *domain.CraftingContext) error {
			called = true
			return nil
		},
	}

	if err := service.ExecutePipeline(ctx, steps); err == nil {
		t.Fatal("expected error, got nil")
	}
	if called {
		t.Fatal("expected pipeline to stop on error")
	}
}

//Catalyst Tests

func TestImbuementCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Normal)

	c := &service.ImbuementCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
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

	c := &service.ReconstructionCatalyst{}
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

	c := &service.ElevatingCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
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

	c := &service.AscendantCatalyst{}

	if err := c.Apply(ctx, domain.Both); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newAffixCount := len(ctx.Item.Prefixes) + len(ctx.Item.Suffixes)

	if newAffixCount <= oldAffixCount {
		t.Fatalf("expected more affixes applied, got %d", newAffixCount)
	}

}

func TestDefiantCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	c := &service.DefiantCatalyst{}
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

	c := &service.LustratingCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
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

	c := &service.LustratingCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
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
