package model

import (
	"time"
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
