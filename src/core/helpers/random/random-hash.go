package random

import (
	"math/rand"
	"strings"

	helpers "github.com/jei-el/vuo.be-backend/src/core/helpers"
)

func NewRandomHash(size uint) string {
	SeedOnce()

	const length int = len(helpers.ALPHABET)
	var sb strings.Builder
	for i := uint(0); i < size; i++ {
		idx := rand.Intn(length)
		r := rune(helpers.ALPHABET[idx])
		sb.WriteRune(r)
	}

	return sb.String()
}
