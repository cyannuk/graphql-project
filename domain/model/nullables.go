package model

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type State uint8

const (
	None   State = 0
	Exists State = 1
	Null   State = 2
)

type NullTime struct {
	Value time.Time
	State State
}

type NullDate struct {
	Value time.Time
	State State
}

type NullInt struct {
	Value int32
	State State
}

type NullBigInt struct {
	Value int64
	State State
}

type NullFloat struct {
	Value float32
	State State
}

type NullDouble struct {
	Value float64
	State State
}

type NullBool struct {
	Value bool
	State State
}

type NullString struct {
	Value string
	State State
}

func (n *NullTime) Set(t time.Time) {
	n.Value = t
	n.State = Exists
}

func (n *NullDate) Set(t time.Time) {
	n.Value = t
	n.State = Exists
}

func (n *NullInt) Set(i int32) {
	n.Value = i
	n.State = Exists
}

func (n *NullBigInt) Set(i int64) {
	n.Value = i
	n.State = Exists
}

func (n *NullFloat) Set(f float32) {
	n.Value = f
	n.State = Exists
}

func (n *NullDouble) Set(f float64) {
	n.Value = f
	n.State = Exists
}

func (n *NullBool) Set(b bool) {
	n.Value = b
	n.State = Exists
}

func (n *NullString) Set(s string) {
	n.Value = s
	n.State = Exists
}

func (n *NullTime) None() {
	n.State = None
}

func (n *NullDate) None() {
	n.State = None
}

func (n *NullInt) None() {
	n.State = None
}

func (n *NullBigInt) None() {
	n.State = None
}

func (n *NullFloat) None() {
	n.State = None
}

func (n *NullDouble) None() {
	n.State = None
}

func (n *NullBool) None() {
	n.State = None
}

func (n *NullString) None() {
	n.State = None
}

func (n *NullTime) Null() {
	n.State = Null
}

func (n *NullDate) Null() {
	n.State = Null
}

func (n *NullInt) Null() {
	n.State = Null
}

func (n *NullBigInt) Null() {
	n.State = Null
}

func (n *NullFloat) Null() {
	n.State = Null
}

func (n *NullDouble) Null() {
	n.State = Null
}

func (n *NullBool) Null() {
	n.State = Null
}

func (n *NullString) Null() {
	n.State = Null
}

func (n *NullTime) SetIfNone(t time.Time) {
	if n.State == None {
		n.Value = t
	}
}

func (n *NullDate) SetIfNone(t time.Time) {
	if n.State == None {
		n.Value = t
	}
}

func (n *NullInt) SetIfNone(i int32) {
	if n.State == None {
		n.Value = i
	}
}

func (n *NullBigInt) SetIfNone(i int64) {
	if n.State == None {
		n.Value = i
	}
}

func (n *NullFloat) SetIfNone(f float32) {
	if n.State == None {
		n.Value = f
	}
}

func (n *NullDouble) SetIfNone(f float64) {
	if n.State == None {
		n.Value = f
	}
}

func (n *NullBool) SetIfNone(b bool) {
	if n.State == None {
		n.Value = b
	}
}

func (n *NullString) SetIfNone(s string) {
	if n.State == None {
		n.Value = s
	}
}

func (n NullTime) TimestampValue() (pgtype.Timestamp, error) {
	var t pgtype.Timestamp
	if n.State == Exists {
		t.Time = n.Value
		t.Valid = true
	} else {
		t.Valid = false
	}
	return t, nil
}

func (n *NullTime) ScanTimestamp(t pgtype.Timestamp) error {
	if !t.Valid {
		n.Null()
	} else {
		n.Set(t.Time)
	}
	return nil
}

func (n NullDate) DateValue() (pgtype.Date, error) {
	var t pgtype.Date
	if n.State == Exists {
		t.Time = n.Value
		t.Valid = true
	} else {
		t.Valid = false
	}
	return t, nil
}

func (n *NullDate) ScanDate(t pgtype.Date) error {
	if !t.Valid {
		n.Null()
	} else {
		n.Set(t.Time)
	}
	return nil
}

func (n NullInt) Int64Value() (pgtype.Int8, error) {
	var i pgtype.Int8
	if n.State == Exists {
		i.Int64 = int64(n.Value)
		i.Valid = true
	} else {
		i.Valid = false
	}
	return i, nil
}

func (n *NullInt) ScanInt64(i pgtype.Int8) error {
	if !i.Valid {
		n.Null()
	} else {
		n.Set(int32(i.Int64))
	}
	return nil
}

func (n NullBigInt) Int64Value() (pgtype.Int8, error) {
	var i pgtype.Int8
	if n.State == Exists {
		i.Int64 = n.Value
		i.Valid = true
	} else {
		i.Valid = false
	}
	return i, nil
}

func (n *NullBigInt) ScanInt64(i pgtype.Int8) error {
	if !i.Valid {
		n.Null()
	} else {
		n.Set(i.Int64)
	}
	return nil
}

func (n NullFloat) Float64Value() (pgtype.Float8, error) {
	var f pgtype.Float8
	if n.State == Exists {
		f.Float64 = float64(n.Value)
		f.Valid = true
	} else {
		f.Valid = false
	}
	return f, nil
}

func (n *NullFloat) ScanFloat64(f pgtype.Float8) error {
	if !f.Valid {
		n.Null()
	} else {
		n.Set(float32(f.Float64))
	}
	return nil
}

func (n NullDouble) Float64Value() (pgtype.Float8, error) {
	var f pgtype.Float8
	if n.State == Exists {
		f.Float64 = n.Value
		f.Valid = true
	} else {
		f.Valid = false
	}
	return f, nil
}

func (n *NullDouble) ScanFloat64(f pgtype.Float8) error {
	if !f.Valid {
		n.Null()
	} else {
		n.Set(f.Float64)
	}
	return nil
}

func (n NullString) TextValue() (pgtype.Text, error) {
	var t pgtype.Text
	if n.State == Exists {
		t.String = n.Value
		t.Valid = true
	} else {
		t.Valid = false
	}
	return t, nil
}

func (n *NullString) ScanText(t pgtype.Text) error {
	if !t.Valid {
		n.Null()
	} else {
		n.Set(t.String)
	}
	return nil
}

func (n NullBool) BoolValue() (pgtype.Bool, error) {
	var b pgtype.Bool
	if n.State == Exists {
		b.Bool = n.Value
		b.Valid = true
	} else {
		b.Valid = false
	}
	return b, nil
}

func (n *NullBool) ScanBool(b pgtype.Bool) error {
	if !b.Valid {
		n.Null()
	} else {
		n.Set(b.Bool)
	}
	return nil
}

func (n *NullBigInt) String() string {
	switch n.State {
	case Exists:
		return strconv.FormatInt(n.Value, 10)
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullInt) String() string {
	switch n.State {
	case Exists:
		return strconv.FormatInt(int64(n.Value), 10)
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullFloat) String() string {
	switch n.State {
	case Exists:
		return strconv.FormatFloat(float64(n.Value), 'f', -1, 32)
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullDouble) String() string {
	switch n.State {
	case Exists:
		return strconv.FormatFloat(n.Value, 'f', -1, 64)
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullBool) String() string {
	switch n.State {
	case Exists:
		return strconv.FormatBool(n.Value)
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullString) String() string {
	switch n.State {
	case Exists:
		return n.Value
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullTime) String() string {
	switch n.State {
	case Exists:
		return n.Value.Format(time.RFC3339Nano)
	case Null:
		return "null"
	default:
		return ""
	}
}

func (n *NullDate) String() string {
	switch n.State {
	case Exists:
		return n.Value.Format(time.DateOnly)
	case Null:
		return "null"
	default:
		return ""
	}
}
