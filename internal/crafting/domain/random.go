package domain

type RandomSource interface {
	Intn(n int) int
}
