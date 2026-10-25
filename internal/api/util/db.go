package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TimestamptzToTimePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}
