package repositories

import (
	shortenBulk "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	"github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type ShortenBulkRepository interface {
	Get(hash string) (*types.RepositoryDTO[shortenBulk.ShortenBulkEntity], error)
	GetOldests(size uint) (map[string]*types.RepositoryDTO[shortenBulk.ShortenBulkEntity], error)
	Post(hash string) error
	IncrementClicks(hash string) error
}
