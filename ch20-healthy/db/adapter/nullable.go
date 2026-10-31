package adapter

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// A nullable datetime column is a pgtype.Timestamptz in the row sqlc
// generates and a *time.Time in the model: nil is NULL. These convert between
// the two.

func timestamptzToTimePtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func timePtrToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}
