package sharelink

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type RandStringGen struct {
	Rand *rand.Rand
}

func NewRandStringGen() *RandStringGen {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &RandStringGen{
		Rand: r,
	}
}

func (g *RandStringGen) generateCode(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[g.Rand.Intn(len(letters))]
	}
	return string(b)
}
