package adapters

import (
	"math/rand"
)

func NewRandomChannel[T any](
	elements []T,
) <-chan T {
	size := len(elements)
	ch := make(chan T, size)
	for ; size > 0; size-- {
		randomIndex := rand.Intn(size)
		lastIndex := size - 1
		if randomIndex != lastIndex {
			elements[randomIndex], elements[lastIndex] =
				elements[lastIndex], elements[randomIndex]
		}
		ch <- elements[lastIndex]
	}

	return ch
}
