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
	// Pre-fill affixes so ClearAffixes has effect
	ctx.Item.Prefixes = []domain.AffixInstance{{}}
	ctx.Item.Suffixes = []domain.AffixInstance{{}}

	c := &service.LustratingCatalyst{}
	if err := c.Apply(ctx, domain.Both); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All affixes should be cleared, rarity inferred
	if len(ctx.Item.Prefixes) == 0 && len(ctx.Item.Suffixes) == 0 {
		// Expected
	} else {
		t.Fatal("expected all affixes cleared")
	}

	// After clearing, rarity should be Normal (no affixes)
	if ctx.Item.Rarity != domain.Normal {
		t.Fatalf("expected Normal rarity after clearing, got %v", ctx.Item.Rarity)
	}
}
