package model

import (
	"encoding/binary"
	"math"
	"time"
)

const unixEpochOffset = 946684800 * 1000000

type value []byte

func (v value) Time() time.Time {
	usec := int64(binary.BigEndian.Uint64(v))
	return time.Unix(unixEpochOffset/1000000+usec/1000000, (unixEpochOffset%1000000*1000)+(usec%1000000*1000)).UTC()
}

func (v value) Date() Date {
	dayOffset := int32(binary.BigEndian.Uint32(v))
	return Date(time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.UTC))
}

func (v value) Bool() bool {
	return v[0] != 0
}

func (v value) Byte() byte {
	return byte(binary.BigEndian.Uint16(v))
}

func (v value) Int() int {
	return int(binary.BigEndian.Uint64(v))
}

func (v value) Int8() int8 {
	return int8(binary.BigEndian.Uint16(v))
}

func (v value) Int16() int16 {
	return int16(binary.BigEndian.Uint16(v))
}

func (v value) Int32() int32 {
	return int32(binary.BigEndian.Uint32(v))
}

func (v value) Int64() int64 {
	return int64(binary.BigEndian.Uint64(v))
}

func (v value) Uint() uint {
	return uint(binary.BigEndian.Uint64(v))
}

func (v value) Uint8() uint8 {
	return uint8(binary.BigEndian.Uint16(v))
}

func (v value) Uint16() uint16 {
	return binary.BigEndian.Uint16(v)
}

func (v value) Uint32() uint32 {
	return binary.BigEndian.Uint32(v)
}

func (v value) Uint64() uint64 {
	return binary.BigEndian.Uint64(v)
}

func (v value) Float32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(v))
}

func (v value) Float64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(v))
}

func (v value) String() string {
	return string(v)
}

func (v value) Role() Role {
	return Role(binary.BigEndian.Uint32(v))
}
