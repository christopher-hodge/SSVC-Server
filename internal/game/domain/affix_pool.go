package domain

import "SSVC-Server/internal/random"

var BaseAffixes = []BaseAffix{
	{
		ID:       "life_flat",
		Name:     "Bull's",
		Type:     Prefix,
		Tags:     []string{"life"},
		BaseMin:  10,
		BaseMax:  13,
		Weight:   1000,
		MinLevel: 1,
	},
	{
		ID:       "life_regen",
		Name:     "Hydra's",
		Type:     Prefix,
		Tags:     []string{"life", "recover"},
		BaseMin:  5,
		BaseMax:  8,
		Weight:   1000,
		MinLevel: 1,
	},
	{
		ID:       "ego_flat",
		Name:     "Sinborn's",
		Type:     Prefix,
		Tags:     []string{"ego"},
		BaseMin:  5,
		BaseMax:  8,
		Weight:   1000,
		MinLevel: 1,
	},
	{
		ID:       "lightning_res",
		Name:     "of Grounding",
		Type:     Suffix,
		Tags:     []string{"lightning", "resistance"},
		BaseMin:  10,
		BaseMax:  13,
		Weight:   1000,
		MinLevel: 1,
	},
	{
		ID:       "cold_res",
		Name:     "of Insulation",
		Type:     Suffix,
		Tags:     []string{"cold", "resistance"},
		BaseMin:  10,
		BaseMax:  13,
		Weight:   1000,
		MinLevel: 1,
	},
	{
		ID:       "fire_res",
		Name:     "of Incombustibility",
		Type:     Suffix,
		Tags:     []string{"fire", "resistance"},
		BaseMin:  10,
		BaseMax:  13,
		Weight:   1000,
		MinLevel: 1,
	},
}

var AffixPool = []AffixDefinition{}

func InitAffixPool(rng random.RNGer) {
	tierCount := 10
	for _, base := range BaseAffixes {
		for tier := 1; tier <= tierCount; tier++ {
			affix := GenerateAffix(base, tier, rng)
			AffixPool = append(AffixPool, affix)
		}
	}
}
