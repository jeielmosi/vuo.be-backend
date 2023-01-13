package adapters

import (
	"errors"
	"math/rand"
)

func SelectRepositories(
	shortenBulkRepositories []*ShortenBulkRepository,
	repositoriesSize int,
) ([]*ShortenBulkRepository, chan ShortenBulkRepository, error) {
	size := len(shortenBulkRepositories)
	fallbacksSize := size - repositoriesSize

	if repositoriesSize > size {
		return shortenBulkRepositories, nil, errors.New("Not enough repositories")
	}

	if repositoriesSize == size {
		return shortenBulkRepositories, nil, errors.New("Not enough fallbacks")
	}

	ch := make(chan *ShortenBulkRepository, fallbacksSize)
	for i := 0; i < fallbacksSize; i++ {
		index := rand.Intn(size - i)
		lastIndex := size - i - 1
		shortenBulkRepositories[index], shortenBulkRepositories[lastIndex] =
			shortenBulkRepositories[lastIndex], shortenBulkRepositories[index]
		ch <- shortenBulkRepositories[lastIndex]
	}

	return shortenBulkRepositories[:size], ch, nil
}
