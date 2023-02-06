package helpers

import (
	"math/rand"
	"sync"

	random "github.com/jei-el/vuo.be-backend/src/core/helpers/random"
)

var once sync.Once

func NewRandomChannel[T any](
	elements *[]T,
) <-chan T {
	random.SeedOnce()

	size := len(*elements)
	permutation := make([]int, size)
	for i := 0; i < size; i++ {
		permutation[i] = i
	}

	ch := make(chan T, size)
	for ; size > 0; size-- {
		randomIndex := rand.Intn(size)
		lastIndex := size - 1
		if randomIndex != lastIndex {
			permutation[randomIndex], permutation[lastIndex] =
				permutation[lastIndex], permutation[randomIndex]
		}
		ch <- (*elements)[permutation[lastIndex]]
	}
	close(ch)

	return ch
}
