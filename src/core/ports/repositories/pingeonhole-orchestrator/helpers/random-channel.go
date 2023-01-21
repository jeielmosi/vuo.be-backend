package ports

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once

func NewRandomChannel[T any](
	elements *[]T,
) <-chan T {
	once.Do(
		func() {
			rand.Seed(time.Now().UTC().UnixNano())
		},
	)

	size := len(*elements)
	permutation := make([]uint, size)
	for i := 0; i < size; i++ {
		permutation[i] = uint(i)
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

	return ch
}
