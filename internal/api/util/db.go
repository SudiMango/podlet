package util

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TimestamptzToTimePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	// 23505 is Postgres's unique_violation SQLSTATE
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
