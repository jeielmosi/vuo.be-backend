package ports

import (
	shortenBulk "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
)

type ShortenBulkRepository interface {
	GetByHash(hash string) (*shortenBulk.ShortenBulkEntity, error)
	GetOldest(size uint) (map[string]*shortenBulk.ShortenBulkEntity, error)
	PostAtHash(hash string) error
	IncrementClicks(hash string) error
}
