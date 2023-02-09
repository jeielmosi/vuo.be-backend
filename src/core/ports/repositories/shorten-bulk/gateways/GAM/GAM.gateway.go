package GAM

import (
	"errors"
	"math/rand"
	"time"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/helpers"
	random "github.com/jei-el/vuo.be-backend/src/core/helpers/random"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
	pigeonhole_shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/helpers/pigeonhole"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type GAMShortenBulkGateway struct {
	pingeonhole *pigeonhole_shorten_bulk.PigeonholeShortenBulkRepository
}

func (g *GAMShortenBulkGateway) Get(hash string) (*entities.ShortenBulkEntity, error) {
	res, err := g.pingeonhole.Get(hash)
	if err != nil {
		return res.Entity, err
	}

	err, _ = g.pingeonhole.IncrementClicks(hash)
	if err != nil {
		return res.Entity, errors.New("Error at count access processing")
	}

	return res.Entity, nil
}

// TODO: verify errors and creation of RepositoryDTO
func (g *GAMShortenBulkGateway) post(hash string, dto *repositories.RepositoryDTO[entities.ShortenBulkEntity]) error {
	err, resCh := g.pingeonhole.Lock(hash)
	defer g.pingeonhole.Unlock(hash)

	backup, err := g.pingeonhole.Get(hash)
	if err != nil {
		return err
	}

	err, resCh = g.pingeonhole.UndoPost(hash, *dto, resCh)
	if err == nil {
		return err
	}

	err, _ = g.pingeonhole.UndoPost(hash, *backup, resCh)

	return err
}

func (g *GAMShortenBulkGateway) postAtNewHash(dto *repositories.RepositoryDTO[entities.ShortenBulkEntity]) (string, error) {
	const TRY_SIZE int = 11

	for t := 0; t < TRY_SIZE; t++ {
		hash := random.NewRandomHash(helpers.HASH_SIZE)
		g.post(hash, dto)
	}

	return "", errors.New("Internal error: Not found empty hash")
}

func (g *GAMShortenBulkGateway) postAtOldHash(
	dto *repositories.RepositoryDTO[entities.ShortenBulkEntity],
) (string, error) {
	const OLDESTS_SIZE uint = 101

	mp, err := g.pingeonhole.GetOldests(OLDESTS_SIZE)
	if err != nil {
		return "", err
	}

	arr := helpers.MapToSlice(mp)

	for size := len(mp); size > 0; size-- {
		idx := rand.Intn(size)
		last := size - 1
		if idx != last {
			arr[idx], arr[last] = arr[last], arr[idx]
		}

		hash := arr[last].Key
		dto := arr[last].Value

		g.post(hash, dto)
	}

	return "", errors.New("Internal error: Hash not found")
}

func (g *GAMShortenBulkGateway) Post(entity entities.ShortenBulkEntity) (string, error) {
	now := time.Now()
	dto := repositories.NewRepositoryDTO(&entity, now, nil)

	res, err := g.postAtNewHash(dto)
	if err == nil {
		return res, err
	}

	return g.postAtOldHash(dto)
}

func NewGAMShortenBulkGateway() (shorten_bulk_gateway.ShortenBulkGateway, error) {
	var repos = &[]*shorten_bulk.ShortenBulkRepository{
		//TODO: Create G.A.M
	}
	pigeonhole, err := pigeonhole_shorten_bulk.NewPigeonholeShortenBulkRepository(repos)
	if err != nil {
		return nil, err
	}

	return &GAMShortenBulkGateway{
		pigeonhole,
	}, nil
}
