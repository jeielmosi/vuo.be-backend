package random

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once

func SeedOnce() {
	once.Do(
		func() {
			rand.Seed(time.Now().UTC().UnixNano())
		},
	)
}
