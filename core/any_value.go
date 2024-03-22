package core

import (
	"encoding/base64"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

type Any struct {
	value any
}

func NewAny(v any) Any {
	return Any{v}
}

func (a Any) IsEmpty() bool {
	return a.value == nil
}

func (a Any) IsDefault() bool {
	switch v := a.value.(type) {
	case string:
		return v == ""
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case bool:
		return v == false
	case []byte:
		return v == nil
	case time.Duration:
		return v == 0
	case netip.Addr:
		return v == netip.Addr{}
	default:
		return true
	}
}

func (a Any) Int() (v int, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(int)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		v = int(n)
	}
	return
}

func (a Any) Int8() (v int8, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(int8)
	if ok {
		return
	}
	n, ok := a.value.(int)
	if ok {
		v = int8(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseInt(s, 10, 16); err != nil {
		return 0, err
	} else {
		return int8(v), nil
	}
}

func (a Any) Int16() (v int16, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(int16)
	if ok {
		return
	}
	n, ok := a.value.(int)
	if ok {
		v = int16(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseInt(s, 10, 16); err != nil {
		return 0, err
	} else {
		return int16(v), nil
	}
}

func (a Any) Int32() (v int32, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(int32)
	if ok {
		return
	}
	n, ok := a.value.(int)
	if ok {
		v = int32(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseInt(s, 10, 32); err != nil {
		return 0, err
	} else {
		return int32(v), nil
	}
}

func (a Any) Int64() (v int64, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(int64)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	} else {
		return int64(v), nil
	}
}

func (a Any) Uint() (v uint, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(uint)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseUint(s, 10, 64); err != nil {
		return 0, err
	} else {
		return uint(v), nil
	}
}

func (a Any) Uint8() (v uint8, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(uint8)
	if ok {
		return
	}
	n, ok := a.value.(uint)
	if ok {
		v = uint8(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseUint(s, 10, 16); err != nil {
		return 0, err
	} else {
		return uint8(v), nil
	}
}

func (a Any) Uint16() (v uint16, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(uint16)
	if ok {
		return
	}
	n, ok := a.value.(uint)
	if ok {
		v = uint16(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseUint(s, 10, 16); err != nil {
		return 0, err
	} else {
		return uint16(v), nil
	}
}

func (a Any) Uint32() (v uint32, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(uint32)
	if ok {
		return
	}
	n, ok := a.value.(uint)
	if ok {
		v = uint32(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseUint(s, 10, 32); err != nil {
		return 0, err
	} else {
		return uint32(v), nil
	}
}

func (a Any) Uint64() (v uint64, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(uint64)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseUint(s, 10, 64); err != nil {
		return 0, err
	} else {
		return uint64(v), nil
	}
}

func (a Any) Float32() (v float32, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(float32)
	if ok {
		return
	}
	n, ok := a.value.(float64)
	if ok {
		v = float32(n)
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseFloat(s, 32); err != nil {
		return 0, err
	} else {
		return float32(v), nil
	}
}

func (a Any) Float64() (v float64, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	v, ok := a.value.(float64)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	if v, err := strconv.ParseFloat(s, 64); err != nil {
		return 0, err
	} else {
		return float64(v), nil
	}
}

func (a Any) Addr() (addr netip.Addr, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	addr, ok := a.value.(netip.Addr)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	return ParseAddress(s)
}

func (a Any) Duration() (t time.Duration, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	t, ok := a.value.(time.Duration)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	return time.ParseDuration(s)
}

func (a Any) Bool() (b bool, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	b, ok := a.value.(bool)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	s = strings.ToLower(s)
	if s == "no" || s == "n" {
		return false, nil
	}
	if s == "yes" || s == "y" {
		return true, nil
	}
	return strconv.ParseBool(s)
}

func (a Any) Bytes() (bytes []byte, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	bytes, ok := a.value.([]byte)
	if ok {
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
		return
	} else if s == "" {
		err = ErrNoValue
		return
	}
	return base64.StdEncoding.DecodeString(s)
}

func (a Any) Str() (s string, err error) {
	if a.value == nil {
		err = ErrNoValue
		return
	}
	s, ok := a.value.(string)
	if !ok {
		err = ErrInvalidValueType
	}
	return
}
