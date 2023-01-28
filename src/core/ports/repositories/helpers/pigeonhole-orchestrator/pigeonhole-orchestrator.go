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

func (o *PigeonholeOrchestrator[T, K]) ExecuteSingleFunc(
	singleFunc func(*T) (*repositories.RepositoryDTO[K], error),
) (*repositories.RepositoryDTO[K], error) {
	if len(*o.repositories) < o.worksSize {
		return nil, errors.New("Internal error: Not enough repositories")
	}
	randomRepositories := helpers.NewRandomChannel(o.repositories)

	var wg sync.WaitGroup
	resultCh := make(chan *repositories.RepositoryDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for repository, ok := <-randomRepositories; ok; repository, ok = <-randomRepositories {
				res, err := singleFunc(repository)
				if err == nil {
					resultCh <- res
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(resultCh)

	if len(resultCh) != o.worksSize {
		return nil, errors.New("Internal error: Not enough successful workers")
	}

	valMp := map[string]*repositories.RepositoryDTO[K]{}
	countMp := map[string]int{}
	for result := range resultCh {
		timestamp := helpers.TimeTo10NanosecondsString(result.UpdatedAt)

		val, ok := valMp[timestamp]
		count := countMp[timestamp]
		if !ok {
			val = result
			count = 0
		}

		valMp[timestamp] = val
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

	return valMp[bestTimestamp], nil
}

// TODO: update by count and timestamp like ExecuteSingleFunc
func (o *PigeonholeOrchestrator[T, K]) ExecuteMultipleFunc(
	multipleFunc func(*T) (map[string]*repositories.RepositoryDTO[K], error),
) (res map[string]*repositories.RepositoryDTO[K], err error) {
	if len(*o.repositories) < o.worksSize {
		return res, errors.New("Internal error: Not enough repositories")
	}
	randomRepositories := helpers.NewRandomChannel(o.repositories)

	var wg sync.WaitGroup
	resultCh := make(chan map[string]*repositories.RepositoryDTO[K], o.worksSize)
	for w := 0; w < o.worksSize; w++ {
		wg.Add(1)
		go func() {
			for repository, ok := <-randomRepositories; ok; repository, ok = <-randomRepositories {
				res, err := multipleFunc(repository)
				if err == nil {
					resultCh <- res
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(resultCh) != o.worksSize {
		return res, errors.New("Internal error: Not enough successful workers")
	}

	valMp := map[string]map[string]*repositories.RepositoryDTO[K]{}
	countMp := map[string]map[string]int{}
	for mp := range resultCh {
		for hash, dto := range mp {
			_, ok := valMp[hash]
			if !ok {
				valMp[hash] = map[string]*repositories.RepositoryDTO[K]{}
				countMp[hash] = map[string]int{}
			}
			timestamp := helpers.TimeTo10NanosecondsString(dto.UpdatedAt)

			val, ok := valMp[hash][timestamp]
			count := countMp[hash][timestamp]
			if !ok {
				val = dto
				count = 0
			}

			valMp[hash][timestamp] = val
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

	return ansMp, nil
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
