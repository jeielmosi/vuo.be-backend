package shorten_bulk

import (
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
)

type ShortenBulkRepository interface {
	Get(hash string) (*types.ShortenBulkDTO, error)
	GetOldests(size uint) (map[string]*types.ShortenBulkDTO, error)
	Post(hash string) error
	IncrementClicks(hash string) error
}
