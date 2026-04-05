package player

const (
	StatLife StatType = iota
	StatLifeRegen
	StatFireRes
	StatColdRes
	StatLightningRes
)

type Modifier struct {
	Type  StatType
	Value float32
	Op    ModifierOp
}

type ModifierOp int

const (
	Add ModifierOp = iota
	Increase
	More
	Override
)
