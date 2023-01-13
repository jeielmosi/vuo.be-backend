package adapters

import (
	"errors"
	"sync"
)

func SingleOperationOrchestrate[T any, K any](
	tools []T,
	worksSize int,
	worker SingleOperationFunction[T, K],
) (res K, err error) {
	toolsSize := len(tools)
	if toolsSize < worksSize {
		return res, errors.New("Internal error: Not enough tools")
	}
	randomTools := NewRandomChannel(tools)

	var wg sync.WaitGroup
	resultCh := make(chan K, worksSize)
	for w := 0; w < worksSize; w++ {
		wg.Add(1)
		go func() {
			for tool, ok := <-randomTools; ok; tool, ok = <-randomTools {
				res1, err1 := worker.work(tool)
				if err1 == nil {
					resultCh <- res1
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(resultCh) != worksSize {
		return res, errors.New("Internal error: Not enough successful workers")
	}

	ans := <-resultCh
	for result := range resultCh {
		if ans < result {
			ans = result
		}
	}
	close(resultCh)

	return ans, nil
}

/*
func MultipleOperationOrchestrate(
	shortenBulkRepositories []*ShortenBulkRepository,
	repositoriesSize int,
) {
	repositories, fallbackCh := SelectRepositories(shortenBulkRepositories, repositoriesSize)
}
*/
