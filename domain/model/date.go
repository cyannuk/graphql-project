package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Date time.Time
type Timestamp time.Time

func (d Date) String() string {
	return (time.Time)(d).Format(time.DateOnly)
}

func (d Date) DateValue() (pgtype.Date, error) {
	return pgtype.Date{Time: time.Time(d), Valid: true}, nil
}

func (t Timestamp) TimestampValue() (pgtype.Timestamp, error) {
	return pgtype.Timestamp{Time: time.Time(t), Valid: true}, nil
}
