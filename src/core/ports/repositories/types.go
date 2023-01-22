package ports

import "time"

type DatabaseDTO[T any] struct {
	Entity    *T
	Locked    bool
	LockedAt  *time.Time
	UpdatedAt time.Time
}

func (lhs *DatabaseDTO[T]) Compare(rhs *DatabaseDTO[T]) int {

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

func NewDatabaseDTO[T any](
	entity *T,
	lockedAt *string,
	updatedAt string,
) (*DatabaseDTO[T], error) {

	lockedAtTime, err := NewTimeFrom10NanosecondsString(lockedAt)
	if err != nil {
		return nil, err
	}

	updatedAtTime, err := NewTimeFrom10NanosecondsString(&updatedAt)
	if err != nil {
		return nil, err
	}

	return &DatabaseDTO[T]{
		Entity:    entity,
		Locked:    (lockedAtTime != nil),
		LockedAt:  lockedAtTime,
		UpdatedAt: *updatedAtTime,
	}, nil
}
