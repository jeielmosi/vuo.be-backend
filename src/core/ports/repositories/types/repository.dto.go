package types

import (
	"time"

	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
)

type RepositoryDTO[T any] struct {
	Entity    *T
	Locked    bool
	LockedAt  *time.Time
	UpdatedAt time.Time
}

func (lhs *RepositoryDTO[T]) Compare(rhs *RepositoryDTO[T]) int {

	if (lhs == nil) && (rhs == nil) {
		return 0
	}
	if (lhs == nil) != (rhs == nil) {
		if lhs == nil {
			return 1
		}
		return -1
	}

	if lhs.UpdatedAt.Before(rhs.UpdatedAt) {
		return -1
	}
	if lhs.UpdatedAt.After(rhs.UpdatedAt) {
		return 1
	}

	if !lhs.Locked && !rhs.Locked {
		return 0
	}
	if lhs.Locked != rhs.Locked {
		if lhs.Locked {
			return 1
		}
		return -1
	}

	if lhs.LockedAt.Before(*rhs.LockedAt) {
		return -1
	}
	if lhs.LockedAt.After(*rhs.LockedAt) {
		return 1
	}

	return 0
}

func NewRepositoryDTO[T any](
	entity *T,
	lockedAt *string,
	updatedAt string,
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
		Locked:    (lockedAtTime != nil),
		LockedAt:  lockedAtTime,
		UpdatedAt: *updatedAtTime,
	}, nil
}
