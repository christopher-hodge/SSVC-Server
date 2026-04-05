package player

type Resource struct {
	Current int
	Max     int
	Regen   float32
}

type Resources struct {
	Life Resource
	Ego  Resource
}

type DamageType int

const (
	Fire DamageType = iota
	Cold
	Lightning
	Psychic
)

type Resistances map[DamageType]float32
type MaxResistances map[DamageType]float32

type Stats struct {
	Base    BaseStats
	Derived DerivedStats
}

type BaseStats struct {
	Life      int
	LifeRegen float32
	Ego       int
	EgoRegen  float32

	CriticalHitChance     float32
	CriticalHitMultiplier float32

	OmniSpeed     float32
	MovementSpeed float32
	CastSpeed     float32
	AttackSpeed   float32

	FireResistance         float32
	MaxFireResistance      float32
	ColdResistance         float32
	MaxColdResistance      float32
	LightningResistance    float32
	MaxLightningResistance float32
	PsychicResistance      float32
	MaxPsychicResistance   float32
}

type DerivedStats struct {
	MaxLife int
	MaxEgo  int

	FinalCritChance     float32
	FinalCritMultiplier float32

	FinalMovementSpeed float32
	FinalCastSpeed     float32
	FinalAttackSpeed   float32

	EffectiveFireRes      float32
	EffectiveColdRes      float32
	EffectiveLightningRes float32
	EffectivePsychicRes   float32
}
