package repository_helpers

import (
	"math/rand"

	random "github.com/jei-el/vuo.be-backend/src/core/helpers/random"
)

func NewRandChIdxs[T any](
	elements *[]T,
) <-chan int {
	random.SeedOnce()

	size := len(*elements)
	permutation := make([]int, size)
	for i := 0; i < size; i++ {
		permutation[i] = i
	}

	ch := make(chan int, size)
	for ; size > 0; size-- {
		randomIndex := rand.Intn(size)
		lastIndex := size - 1
		if randomIndex != lastIndex {
			permutation[randomIndex], permutation[lastIndex] =
				permutation[lastIndex], permutation[randomIndex]
		}
		ch <- permutation[lastIndex]
	}
	close(ch)

	return ch
}
