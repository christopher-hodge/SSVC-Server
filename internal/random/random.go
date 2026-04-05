package random

import (
	"math/rand"
	"time"
)

type RNGer interface {
	Intn(n int) int
}

type RNGerFloat interface {
	Floatn(n float64) float64
}

type RNG struct {
	r    *rand.Rand
	seed int64
}

func New(seed int64) *RNG {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &RNG{
		r:    rand.New(rand.NewSource(seed)),
		seed: seed,
	}
}

func (rng *RNG) Intn(n int) int {
	return rng.r.Intn(n)
}

func (rng *RNG) Seed() int64 {
	return rng.seed
}
