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

func (lhs *RepositoryDTO[T]) IsOlderThan(rhs *RepositoryDTO[T]) bool {

	if (lhs == nil) && (rhs == nil) {
		return false
	}
	if (lhs == nil) != (rhs == nil) {
		if lhs == nil {
			return true
		}
		return false
	}

	if lhs.UpdatedAt.Before(rhs.UpdatedAt) {
		return true
	}
	if lhs.UpdatedAt.After(rhs.UpdatedAt) {
		return false
	}

	if !lhs.IsLocked() && !rhs.IsLocked() {
		return false
	}
	if lhs.IsLocked() != rhs.IsLocked() {
		if lhs.IsLocked() {
			return true
		}
		return false
	}

	if lhs.LockedAt.Before(*rhs.LockedAt) {
		return true
	}
	if lhs.LockedAt.After(*rhs.LockedAt) {
		return false
	}

	return false
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
