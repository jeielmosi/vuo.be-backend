package helpers

import (
	"time"
)

const _10NanosecondsPattern = "2006-01-02T15:04:05.99999999Z"

func NewTimeFrom10NanosecondsString(timestamp *string) (*time.Time, error) {
	if timestamp == nil {
		return nil, nil
	}
	res, err := time.Parse(_10NanosecondsPattern, *timestamp)

	return &res, err
}

func TimeTo10NanosecondsString(t time.Time) string {
	return t.UTC().Format(_10NanosecondsPattern)
}

func Now10NanosecondsString() string {
	return time.Now().UTC().Format(_10NanosecondsPattern)
}
