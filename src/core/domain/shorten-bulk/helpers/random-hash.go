package helpers

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

var once sync.Once

func NewRandomHash(size uint) string {
	once.Do(
		func() {
			rand.Seed(time.Now().UTC().UnixNano())
		},
	)

	//TODO: verify if could increase alphabet

	const alphabet string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789-abcdefghijklmnopqrstuvwxyz"
	const length int = len(alphabet)

	var sb strings.Builder
	for i := uint(0); i < size; i++ {
		idx := rand.Intn(length)
		r := rune(alphabet[idx])
		sb.WriteRune(r)
	}

	return sb.String()
}
