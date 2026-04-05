package modifier

type StatType int

const (
	StatLife StatType = iota
	StatLifeRegen
	StatFireRes
	StatColdRes
	StatLightningRes
)

type Modifier struct {
	Type  StatType
	Value []float64
	Op    ModifierOp
}

type ModifierOp int

const (
	Add ModifierOp = iota
	Increase
	More
	Override
)
