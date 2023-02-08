package usecases

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
)

func GetShortenBulkEntity(gateway *shorten_bulk_gateway.ShortenBulkGateway, hash string) (
	*entities.ShortenBulkEntity,
	error,
) {
	return (*gateway).Get(hash)
}
