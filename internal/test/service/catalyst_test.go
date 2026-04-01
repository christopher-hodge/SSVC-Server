package service

import (
	"errors"
	"testing"

	"SSVC-Server/internal/crafting/domain"
	"SSVC-Server/internal/crafting/service"
	"SSVC-Server/internal/random"
)

// test ApplyAffix function
func testApplyAffix(ctx *domain.CraftingContext, t domain.AffixType) error {
	ctx.Item.Prefixes = append(ctx.Item.Prefixes, domain.AffixInstance{})
	ctx.Item.Suffixes = append(ctx.Item.Suffixes, domain.AffixInstance{})
	return nil
}

//
// --- Helpers ---
//

var rarityNames = map[domain.Rarity]string{
	domain.Normal: "Normal",
	domain.Magic:  "Magic",
	domain.Rare:   "Rare",
}

func newTestItem(rarity domain.Rarity) *domain.Item {
	return &domain.Item{
		Rarity:   rarity,
		Prefixes: []domain.AffixInstance{},
		Suffixes: []domain.AffixInstance{},
	}
}

func newTestContext(rarity domain.Rarity) *domain.CraftingContext {
	return &domain.CraftingContext{
		Item: newTestItem(rarity),
		RNG:  *random.New(42),
	}
}

//
// --- Pipeline Tests ---
//

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

//
// --- Catalyst Tests ---
//

func TestImbuementCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Normal)

	c := &service.ImbuementCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Item.Rarity != domain.Magic {
		t.Fatalf("expected Magic rarity, got %v", ctx.Item.Rarity)
	}
	if len(ctx.Item.Prefixes) != 1 || len(ctx.Item.Suffixes) != 1 {
		t.Fatal("expected 1 prefix and 1 suffix applied")
	}
}

func TestElevatingCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Magic)

	c := &service.ElevatingCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ctx.Item.Rarity != domain.Rare {
		t.Fatalf("expected Rare rarity, got %v", ctx.Item.Rarity)
	}
	if len(ctx.Item.Prefixes) != 1 || len(ctx.Item.Suffixes) != 1 {
		t.Fatal("expected 1 prefix and 1 suffix applied")
	}
}

func TestAscendantCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	c := &service.AscendantCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ctx.Item.Prefixes) != 1 || len(ctx.Item.Suffixes) != 1 {
		t.Fatal("expected 1 prefix and 1 suffix applied")
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
}

func TestLustratingCatalyst(t *testing.T) {
	ctx := newTestContext(domain.Rare)

	ctx.Item.Prefixes = []domain.AffixInstance{{}}
	ctx.Item.Suffixes = []domain.AffixInstance{{}}

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

	// Pre-fill affixes
	ctx.Item.Prefixes = []domain.AffixInstance{{DefID: "test", Value: 1, Tags: []string{"test"}}, {DefID: "test2", Value: 1, Tags: []string{"test2"}}}
	ctx.Item.Suffixes = []domain.AffixInstance{{DefID: "lock_prefixes", Value: 1, Tags: []string{"lock_prefixes"}}, {}}

	c := &service.LustratingCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Prefixes should remain because they are locked
	if len(ctx.Item.Prefixes) != 2 {
		t.Fatalf("expected prefixes to remain, got %d", len(ctx.Item.Prefixes))
	}

	// Suffixes should have been cleared
	if len(ctx.Item.Suffixes) != 0 {
		t.Fatalf("expected suffixes cleared, got %d", len(ctx.Item.Suffixes))
	}

	// Rarity should be inferred from remaining mods
	if ctx.Item.Rarity != domain.Rare { // depending on limits
		t.Fatalf("expected Rare rarity after clearing, got %s", rarityNames[ctx.Item.Rarity])
	}
}
