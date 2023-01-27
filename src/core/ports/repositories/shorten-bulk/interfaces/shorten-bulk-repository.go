package shorten_bulk

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type ShortenBulkRepository interface {
	Get(hash string) (*repositories.RepositoryDTO[entities.ShortenBulkEntity], error)
	GetOldests(size uint) (map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity], error)
	Post(hash string, shorten_bulk repositories.RepositoryDTO[entities.ShortenBulkEntity]) error
	IncrementClicks(hash string) error
}
