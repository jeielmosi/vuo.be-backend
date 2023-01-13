package adapters

import (
	"domain"
	"time"
)

type LockDTO struct {
	Locked   bool
	LockedAt time.Time
}

func NewLock(
	lock bool,
	lockedAt string,
) *LockDTO {
	return &LockDTO{
		Locked:   lock,
		LockedAt: NewTimeFrom10NanosecondsString(updatedAt),
	}
}

type ShortenBulkDTO struct {
	ShortenBulk *domain.ShortenBulkEntity
	Lock        *LockDTO
	UpdatedAt   time.Time
}

func NewShortenBulkDTO(
	shortenBulk *domain.ShortenBulkEntity,
	lock *LockDTO,
	updatedAt string,
) *ShortenBulkDTO {
	return &ShortenBulkDTO{
		ShortenBulk: shortenBulk,
		Lock:        lock,
		UpdatedAt:   NewTimeFrom10NanosecondsString(updatedAt),
	}
}

func (lhs *ShortenBulkDTO) Compare(rhs *ShortenBulkDTO) int {
	if lhs.UpdatedAt.Before(rhs.UpdatedAt) {
		return -1
	}

	if lhs.UpdatedAt.After(rhs.UpdatedAt) {
		return 1
	}

	return 0
}
