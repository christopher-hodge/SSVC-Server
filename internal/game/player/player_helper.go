package player

func CalculateStat(base float32, mods []Modifier) float32 {
	flat := float32(0)
	increase := float32(0)
	more := float32(1)

	for _, m := range mods {
		switch m.Op {
		case Add:
			flat += m.Value

		case Increase:
			increase += m.Value

		case More:
			more *= (1 + m.Value)

		case Override:
			return m.Value
		}
	}

	return (base + flat) * (1 + increase) * more
}
