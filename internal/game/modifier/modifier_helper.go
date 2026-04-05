package modifier

func CalculateStat(base float64, mods []Modifier, valueIndex int) float64 {
	flat := float64(0)
	increase := float64(0)
	more := float64(1)

	for _, m := range mods {
		switch m.Op {
		case Add:
			flat += m.Value[valueIndex]

		case Increase:
			increase += m.Value[valueIndex]

		case More:
			more *= (1 + m.Value[valueIndex])

		case Override:
			return m.Value[valueIndex]
		}
	}

	return (base + flat) * (1 + increase) * more
}

func BuildModifiers(id string, value []float64) []Modifier {
	switch id {

	case "life_flat":
		return []Modifier{
			{
				Type:  StatLife,
				Value: value,
				Op:    Add,
			},
		}

	case "life_regen":
		return []Modifier{
			{
				Type:  StatLifeRegen,
				Value: value,
				Op:    Add,
			},
		}

	case "fire_res":
		return []Modifier{
			{
				Type:  StatFireRes,
				Value: value,
				Op:    Add,
			},
		}

	case "cold_res":
		return []Modifier{
			{
				Type:  StatColdRes,
				Value: value,
				Op:    Add,
			},
		}

	case "lightning_res":
		return []Modifier{
			{
				Type:  StatLightningRes,
				Value: value,
				Op:    Add,
			},
		}
	}

	return nil
}
