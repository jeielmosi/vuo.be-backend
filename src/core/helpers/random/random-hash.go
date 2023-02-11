package random

import (
	"math/rand"
	"strings"

	helpers "github.com/jei-el/vuo.be-backend/src/core/helpers"
)

func NewRandomHash(size uint) string {
	if size == 0 {
		return ""
	}

	SeedOnce()

	var sb strings.Builder

	idx := rand.Intn(len(helpers.FIRST_CHAR_ALPHABET))
	r := rune(helpers.FIRST_CHAR_ALPHABET[idx])
	sb.WriteRune(r)

	for i := uint(1); i < size; i++ {
		idx = rand.Intn(len(helpers.ALPHABET))
		r = rune(helpers.ALPHABET[idx])
		sb.WriteRune(r)
	}

	return sb.String()
}
