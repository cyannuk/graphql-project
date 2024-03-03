package model

import (
	"time"
)

type UserSession struct {
	Name   string
	Email  string
	UserId int64
	Admin  bool
	Exp    time.Time
}
