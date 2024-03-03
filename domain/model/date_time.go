package model

import (
	"bytes"
	"time"
)

type DateTime time.Time

type Date time.Time

func (dt *Date) UnmarshalJSON(b []byte) error {
	t, err := time.ParseInLocation(`"2006-01-02"`, string(b), time.UTC)
	if err != nil {
		return err
	}
	*dt = Date(t)
	return nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	var layout string
	if bytes.HasSuffix(b, []byte(`Z"`)) {
		layout = `"2006-01-02T15:04:05.999999Z"`
	} else {
		layout = `"2006-01-02T15:04:05.999999"`
	}
	t, err := time.ParseInLocation(layout, string(b), time.UTC)
	if err != nil {
		return err
	}
	*dt = DateTime(t)
	return nil
}

func (dt *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(*dt)
	return []byte(t.Format("2006-01-02")), nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(*dt)
	return []byte(t.Format("2006-01-02T15:04:05.999Z")), nil
}
