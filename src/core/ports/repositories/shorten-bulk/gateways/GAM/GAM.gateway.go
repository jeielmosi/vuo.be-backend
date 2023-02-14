package GAM

import (
	"errors"
	"math/rand"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/helpers"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	firestore_shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/firestore"
	pigeonhole_shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/adapters/pigeonhole"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type GAMShortenBulkGateway struct {
	repository shorten_bulk.ShortenBulkRepository
}

func (g *GAMShortenBulkGateway) Get(hash string) (*entities.ShortenBulkEntity, error) {
	res, err := g.repository.Get(hash)
	if err != nil {
		return res.Entity, err
	}

	err = g.repository.IncrementClicks(hash)
	if err != nil {
		return res.Entity, errors.New("Error at count access processing")
	}

	return res.Entity, nil
}

func (g *GAMShortenBulkGateway) post(
	hash string,
	shortenBulk *entities.ShortenBulkEntity,
	stopFunc func(*repositories.RepositoryDTO[entities.ShortenBulkEntity]) error,
) error {
	err := g.repository.Lock(hash)
	if err != nil {
		return err
	}
	defer g.repository.Unlock(hash)

	backup, err := g.repository.Get(hash)
	if err != nil {
		return err
	}

	err = stopFunc(backup)
	if err != nil {
		return err
	}

	dto := repositories.NewRepositoryDTO(shortenBulk, true)
	err = g.repository.Post(hash, *dto)
	if err == nil {
		return err
	}

	dto = backup.Update()
	return g.repository.Post(hash, *dto)
}

func (g *GAMShortenBulkGateway) postAtNewHash(shortenBulk *entities.ShortenBulkEntity) (string, error) {
	const TRY_SIZE int = 11

	stopFunc := func(dto *repositories.RepositoryDTO[entities.ShortenBulkEntity]) error {
		if dto != nil {
			return errors.New("Hash: is used")
		}
		return nil
	}

	for t := 0; t < TRY_SIZE; t++ {
		hash := helpers.NewRandomHash(helpers.HASH_SIZE)
		err := g.post(hash, shortenBulk, stopFunc)
		if err == nil {
			return hash, err
		}
	}

	return "", errors.New("Internal error: Not found empty hash")
}

func (g *GAMShortenBulkGateway) postAtOldHash(
	shortenBulk *entities.ShortenBulkEntity,
) (string, error) {
	const OLDESTS_SIZE int = 101

	mp, err := g.repository.GetOldest(OLDESTS_SIZE)
	if err != nil {
		return "", err
	}

	keys := helpers.GetKeys(mp)

	lastTimestamp := ""
	for _, key := range keys {
		timestamp := repository_helpers.TimeTo10NanosecondsString(mp[key].CreatedAt)
		if lastTimestamp < timestamp {
			lastTimestamp = timestamp
		}
	}

	stopFunc := func(dto *repositories.RepositoryDTO[entities.ShortenBulkEntity]) error {
		if dto == nil {
			return errors.New("Empty hash")
		}

		if dto.Locked {
			return errors.New("Locked hash")
		}

		timestamp := repository_helpers.TimeTo10NanosecondsString(dto.CreatedAt)
		if timestamp > lastTimestamp {
			return errors.New("Element is not old")
		}
		return nil
	}

	for size := len(mp); size > 0; size-- {
		idx := rand.Intn(size)
		last := size - 1
		if idx != last {
			keys[idx], keys[last] = keys[last], keys[idx]
		}

		hash := keys[last]
		dto := mp[hash]

		err = g.post(hash, dto.Entity, stopFunc)
		if err == nil {
			return hash, err
		}
	}

	return "", errors.New("Internal error: Hash not found")
}

func (g *GAMShortenBulkGateway) Post(shortenBulk entities.ShortenBulkEntity) (string, error) {
	res, err := g.postAtNewHash(&shortenBulk)
	if err == nil {
		return res, err
	}

	return g.postAtOldHash(&shortenBulk)
}

func NewGAMShortenBulkGateway(envName string) (shorten_bulk_gateway.ShortenBulkGateway, error) {
	firestore := firestore_shorten_bulk.NewShortenBulkFirestore(envName)
	var repos = &[]*shorten_bulk.ShortenBulkRepository{
		&firestore,
		//TODO: Create A.M
	}
	pigeonhole, err := pigeonhole_shorten_bulk.NewPigeonholeShortenBulkRepository(repos)
	if err != nil {
		return nil, err
	}

	return &GAMShortenBulkGateway{
		repository: pigeonhole,
	}, nil
}
