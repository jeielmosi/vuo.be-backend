package shorten_bulk_gateway

import entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"

type ShortenBulkRepository interface {
	Get(hash string) (*entities.ShortenBulkEntity, error)
	GetOldests(size uint) (map[string]*entities.ShortenBulkEntity, error)
	Post(hash string, shorten_bulk entities.ShortenBulkEntity) error
	IncrementClicks(hash string) error
}
