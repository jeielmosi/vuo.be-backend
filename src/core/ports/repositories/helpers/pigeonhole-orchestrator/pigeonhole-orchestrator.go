package helpers

import (
	"errors"
	"sync"

	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type PigeonholeOrchestrator[T any, K any] struct {
	worksSize    int
	repositories *[]*T
}

func (o *PigeonholeOrchestrator[T, K]) executeSingleFnWithCh(
	singleFn func(*T) (*repositories.RepositoryDTO[K], error),
	idxsCh <-chan int,
	worksSize int,
) (*repositories.RepositoryDTO[K], error, <-chan int) {
	if len(idxsCh) < worksSize {
		return nil, errors.New("Internal error: Not enough repositories"), nil
	}

	var wg sync.WaitGroup
	resCh := make(chan *repositories.RepositoryDTO[K], o.worksSize)
	resIdxsCh := make(chan int, o.worksSize)
	defer close(resIdxsCh)

	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for idx, ok := <-idxsCh; ok; idx, ok = <-idxsCh {
				repo := (*o.repositories)[idx]
				res, err := singleFn(repo)
				if err == nil {
					resCh <- res
					resIdxsCh <- idx
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(resCh)

	if len(resCh) != o.worksSize {
		return nil, errors.New("Internal error: Not enough successful workers"), resIdxsCh
	}

	valMp := map[string]*repositories.RepositoryDTO[K]{}
	countMp := map[string]int{}
	for result := range resCh {
		timestamp := helpers.TimeTo10NanosecondsString(result.UpdatedAt)

		count, ok := countMp[timestamp]
		if !ok {
			valMp[timestamp] = result
			count = 0
		}
		countMp[timestamp] = count + 1
	}

	bestTimestamp := ""
	bestCount := 0
	for timestamp, count := range countMp {
		if (bestCount < count) || (bestCount == count && bestTimestamp < timestamp) {
			bestCount = count
			bestTimestamp = timestamp
		}
	}

	return valMp[bestTimestamp], nil, resIdxsCh
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteSingleFnWithCh(
	singleFn func(*T) (*repositories.RepositoryDTO[K], error),
	idxsCh <-chan int,
) (*repositories.RepositoryDTO[K], error, <-chan int) {
	return o.executeSingleFnWithCh(singleFn, idxsCh, len(idxsCh))
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteSingleFn(
	singleFn func(*T) (*repositories.RepositoryDTO[K], error),
) (*repositories.RepositoryDTO[K], error, <-chan int) {
	idxsCh := helpers.NewRandChIdxs(o.repositories)
	return o.ExecuteSingleFnWithCh(singleFn, idxsCh)
}

func (o *PigeonholeOrchestrator[T, K]) executeMultipleFnWithCh(
	multipleFn func(*T) (map[string]*repositories.RepositoryDTO[K], error),
	idxsCh <-chan int,
	worksSize int,
) (map[string]*repositories.RepositoryDTO[K], error, <-chan int) {
	res := map[string]*repositories.RepositoryDTO[K]{}
	resIdxsCh := make(chan int, worksSize)
	defer close(resIdxsCh)

	if len(*o.repositories) < worksSize {
		return res, errors.New("Internal error: Not enough repositories"), resIdxsCh
	}

	var wg sync.WaitGroup
	resCh := make(chan map[string]*repositories.RepositoryDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for idx, ok := <-idxsCh; ok; idx, ok = <-idxsCh {
				repo := (*o.repositories)[idx]
				res, err := multipleFn(repo)
				if err == nil {
					resCh <- res
					resIdxsCh <- idx
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(resCh)

	if len(resCh) != o.worksSize {
		return res, errors.New("Internal error: Not enough successful workers"), resIdxsCh
	}

	valMp := map[string]map[string]*repositories.RepositoryDTO[K]{}
	countMp := map[string]map[string]int{}
	for mp := range resCh {
		for hash, dto := range mp {
			_, ok := valMp[hash]
			if !ok {
				valMp[hash] = map[string]*repositories.RepositoryDTO[K]{}
				countMp[hash] = map[string]int{}
			}
			timestamp := helpers.TimeTo10NanosecondsString(dto.UpdatedAt)

			count, _ := countMp[hash][timestamp]
			if !ok {
				valMp[hash][timestamp] = dto
				count = 0
			}
			countMp[hash][timestamp] = count + 1
		}
	}

	ansMp := map[string]*repositories.RepositoryDTO[K]{}
	for hash, mp := range countMp {
		countSum := 0
		bestTimestamp := ""
		bestCount := 0
		for timestamp, count := range mp {
			if (bestCount < count) || (bestCount == count && bestTimestamp < timestamp) {
				bestCount = count
				bestTimestamp = timestamp
			}
			countSum += count
		}
		if countSum <= o.worksSize {
			ansMp[hash] = valMp[hash][bestTimestamp]
		}
	}

	return ansMp, nil, resIdxsCh
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteMultipleFn(
	multipleFn func(*T) (map[string]*repositories.RepositoryDTO[K], error),
) (map[string]*repositories.RepositoryDTO[K], error, <-chan int) {
	idxsCh := helpers.NewRandChIdxs(o.repositories)
	return o.executeMultipleFnWithCh(multipleFn, idxsCh, o.worksSize)
}

func (o *PigeonholeOrchestrator[T, K]) ExecuteMultipleFnWithChannel(
	multipleFn func(*T) (map[string]*repositories.RepositoryDTO[K], error),
	idxsCh <-chan int,
) (map[string]*repositories.RepositoryDTO[K], error, <-chan int) {
	return o.executeMultipleFnWithCh(multipleFn, idxsCh, len(idxsCh))
}

func NewPigeonholeOrchestrator[T any, K any](
	repositories *[]*T,
) (*PigeonholeOrchestrator[T, K], error) {
	if repositories == nil {
		return nil, errors.New("Internal error: Repositories is a nil pointer")
	}

	size := len(*repositories)
	if size == 0 {
		return nil, errors.New("Internal error: Repositories array is empty")
	}

	return &PigeonholeOrchestrator[T, K]{
		worksSize:    size/2 + 1,
		repositories: repositories,
	}, nil
}
