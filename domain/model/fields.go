package model

import (
	"strconv"

	gotils "github.com/savsgio/gotils/strconv"
)

type fields struct {
	names        []byte
	placeholders []byte
	args         []any
}

func (f *fields) addField(name string, value any) {
	l := len(f.args)
	if l > 0 {
		f.names = append(f.names, ',')
		f.placeholders = append(f.placeholders, ',')
	}
	f.names = append(f.names, '"')
	f.names = append(f.names, name...)
	f.names = append(f.names, '"')
	f.placeholders = append(f.placeholders, '$')
	f.placeholders = append(f.placeholders, strconv.FormatInt(int64(l+1), 10)...)
	f.args = append(f.args, value)
}

func (f *fields) get() (string, string, []any) {
	return gotils.B2S(f.names), gotils.B2S(f.placeholders), f.args
}
