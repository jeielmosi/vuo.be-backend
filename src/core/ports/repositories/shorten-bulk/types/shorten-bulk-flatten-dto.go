package types

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

const (
	CreatedAtField = "created_at"
	UpdatedAtField = "updated_at"
	URLField       = "url"
	LockedField    = "locked"
	ClicksField    = "clicks"
)

type ShortenBulkFlattenDTO = map[string]interface{}

func NewShortenBulkFlattenDTO(
	dto repositories.RepositoryDTO[entities.ShortenBulkEntity],
) ShortenBulkFlattenDTO {
	ans := ShortenBulkFlattenDTO{}

	ans[CreatedAtField] = helpers.TimeTo10NanosecondsString(dto.CreatedAt)
	ans[UpdatedAtField] = helpers.TimeTo10NanosecondsString(dto.UpdatedAt)
	ans[URLField] = dto.Entity.URL
	ans[LockedField] = dto.Locked
	ans[ClicksField] = dto.Entity.Clicks

	return ans
}

func ToRepositoryDTO(flatten ShortenBulkFlattenDTO) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	createdAtStr := flatten[CreatedAtField].(string)
	createdAt, err := repository_helpers.NewTimeFrom10NanosecondsString(&createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAtStr := flatten[UpdatedAtField].(string)
	updatedAt, err := repository_helpers.NewTimeFrom10NanosecondsString(&updatedAtStr)
	if err != nil {
		return nil, err
	}

	dto := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity: entities.NewShortenBulkEntity(
			flatten[URLField].(string),
			flatten[ClicksField].(int64),
		),
		CreatedAt: *createdAt,
		Locked:    flatten[LockedField].(bool),
		UpdatedAt: *updatedAt,
	}

	return dto, err
}
