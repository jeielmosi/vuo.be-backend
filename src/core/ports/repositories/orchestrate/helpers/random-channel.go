package adapters

import (
	"math/rand"
)

func NewRandomChannel[T any](
	elements []T,
) chan T {
	size := len(elements)

	ch := make(chan T, size)
	for i := 0; i < size; i++ {
		randomIndex := rand.Intn(size - i)
		lastIndex := size - i - 1
		elements[randomIndex], elements[lastIndex] =
			elements[lastIndex], elements[randomIndex]
		ch <- elements[lastIndex]
	}

	return ch
}
