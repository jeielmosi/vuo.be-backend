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
	if len(tools) < worksSize {
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

	return ans.Entity, nil
}

type ValueCount[T any] struct {
	Value T
	Count uint
}

func MultipleOperationOrchestrate[T any, K any](
	tools []T,
	worksSize int,
	worker MultipleOperationFunction[T, K],
) (res map[string]K, err error) {
	if len(tools) < worksSize {
		return res, errors.New("Internal error: Not enough tools")
	}
	randomTools := NewRandomChannel(tools)

	var wg sync.WaitGroup
	resultCh := make(chan map[string]K, worksSize)
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

	valueCountMap := map[string]ValueCount[PigeonholeDTO[K]]{}
	for resultMap := range resultCh {
		for key, newValue := range resultMap {
			if _, ok := valueCountMap[key]; !ok {
				valueCountMap[key] = ValueCount[PigeonholeDTO[K]]{
					Value: newValue,
					Count: 0,
				}
			}
			if valueCountMap[key].Value < newValue {
				valueCountMap[key].Value = newValue
			}
			valueCountMap[key].Count++
		}
	}

	ans := map[string]K{}
	for key, valueCount := range valueCountMap {
		if valueCount.Count == worksSize {
			ans[key] = valueCount.Value.Entity
		}
	}

	return ans, nil
}
