package adapters

import (
	"errors"
	"sync"
)

func SingleOperationOrchestrate[K any](
	shortenBulkRepositories []*ShortenBulkRepository,
	repositoriesSize int,
	f SingleOperationFunction,
) (K, error) {
	repositories, fallbackCh, err := SelectRepositories(shortenBulkRepositories, repositoriesSize)
	if err != nil {
		return K{}, err
	}

	var wg sync.WaitGroup
	resultsCh := make(chan K, repositoriesSize)
	for _, repository := range repositories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.work(repository, fallbackCh, resultsCh)
		}()
	}
	wg.Wait()

	if len(resultsCh) != repositoriesSize {
		return K{}, errors.New("Internal error")
	}

	//TODO: orchestrate the return at resultsCh
}

/*
func MultipleOperationOrchestrate(
	shortenBulkRepositories []*ShortenBulkRepository,
	repositoriesSize int,
) {
	repositories, fallbackCh := SelectRepositories(shortenBulkRepositories, repositoriesSize)
}
*/
