package domain

var AffixPool = []AffixDefinition{
	{
		ID:       "life_flat_t1",
		Name:     "Bull's",
		Type:     Prefix,
		Tags:     []string{"life"},
		MinValue: 70,
		MaxValue: 89,
		Weight:   100,
		MinLevel: 70,
	},
	{
		ID:       "fire_res_t1",
		Name:     "of Incombustibility",
		Type:     Suffix,
		Tags:     []string{"fire", "resistance"},
		MinValue: 41,
		MaxValue: 45,
		Weight:   120,
		MinLevel: 60,
	},
}
