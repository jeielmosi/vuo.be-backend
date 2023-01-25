package types

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type ShortenBulkDTO = types.RepositoryDTO[domain.ShortenBulkEntity]
