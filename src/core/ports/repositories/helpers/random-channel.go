package repository_helpers

import (
	"math/rand"
)

func NewRandChanIdxs(size uint) <-chan int {
	perm := rand.Perm(int(size))

	ch := make(chan int, size)
	for _, val := range perm {
		ch <- val
	}
	close(ch)

	return ch
}
