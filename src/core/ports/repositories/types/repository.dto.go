package repositories

import (
	"time"
)

type RepositoryDTO[T any] struct {
	Entity    *T
	LockedAt  *time.Time
	UpdatedAt time.Time
}

func (r *RepositoryDTO[T]) IsLocked() bool {
	return r.LockedAt != nil
}

func NewRepositoryDTO[T any](
	entity *T,
	updatedAt time.Time,
	lockedAt *time.Time,
) *RepositoryDTO[T] {

	return &RepositoryDTO[T]{
		Entity:    entity,
		LockedAt:  lockedAt,
		UpdatedAt: updatedAt,
	}
}

/*
func NewRepositoryDTO[T any](
	entity *T,
	lockedAt *time.Time,
	updatedAt *time.Time,
) (*RepositoryDTO[T], error) {

	lockedAtTime, err := helpers.NewTimeFrom10NanosecondsString(lockedAt)
	if err != nil {
		return nil, err
	}

	updatedAtTime, err := helpers.NewTimeFrom10NanosecondsString(&updatedAt)
	if err != nil {
		return nil, err
	}

	return &RepositoryDTO[T]{
		Entity:    entity,
		LockedAt:  lockedAtTime,
		UpdatedAt: *updatedAtTime,
	}, nil
}
*/
