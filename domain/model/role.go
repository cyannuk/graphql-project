package model

type Role int32

const (
	RoleRefresh Role = -2
	RoleAnon    Role = -1
	RoleUser    Role = 0
	RoleAdmin   Role = 10
)
