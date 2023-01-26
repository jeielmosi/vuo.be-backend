package shorten_bulk

import (
	domain "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	pigeonhole_orchestrator "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers/pigeonhole-orchestrator"
)

type MultipleOperation = pigeonhole_orchestrator.MultipleOperation[ShortenBulkRepository, domain.ShortenBulkEntity]
