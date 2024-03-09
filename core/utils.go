package core

import (
	"github.com/savsgio/gotils/strconv"
)

func TrimQuotes(str string) string {
	if str == "" {
		return str
	}
	bytes := strconv.S2B(str)
	first := bytes[0]
	l := len(bytes) - 1
	last := bytes[l]
	if (first == '\'' && last == '\'') || (first == '"' && last == '"') || (first == '`' && last == '`') {
		return strconv.B2S(bytes[1:l])
	}
	return str
}

func StartWith(str string, ch byte) bool {
	if str == "" {
		return false
	}
	for _, b := range strconv.S2B(str) {
		if b == ch {
			return true
		}
		if b != ' ' && b != '\t' {
			return false
		}
	}
	return false
}

// ASCII only
func IsUpperCase(str string, i int) bool {
	if str == "" {
		return false
	}
	bytes := strconv.S2B(str)
	if i < len(bytes) {
		b := bytes[i]
		return b >= 'A' && b <= 'Z'
	}
	return false
}

// ASCII only
func Capitalize(str string) string {
	if str == "" {
		return str
	}
	bytes := strconv.S2B(str)
	var b byte
	var i int
	for i, b = range bytes {
		if b != '_' {
			break
		}
	}
	if b >= 'a' && b <= 'z' {
		buffer := make([]byte, 0, len(bytes))
		buffer = append(buffer, b-0x20)
		buffer = append(buffer, bytes[i+1:]...)
		return strconv.B2S(buffer)
	} else {
		if i == 0 {
			return str
		} else {
			return strconv.B2S(bytes[i:])
		}
	}
}

// ASCII only
func Uncapitalize(str string) string {
	if str == "" {
		return str
	}
	bytes := strconv.S2B(str)
	var b byte
	var i int
	for i, b = range bytes {
		if b != '_' {
			break
		}
	}
	if b >= 'A' && b <= 'Z' {
		buffer := make([]byte, 0, len(bytes))
		buffer = append(buffer, b+0x20)
		buffer = append(buffer, bytes[i+1:]...)
		return strconv.B2S(buffer)
	} else {
		if i == 0 {
			return str
		} else {
			return strconv.B2S(bytes[i:])
		}
	}
}

// ASCII only
func Plural(str string) string {
	if str == "" {
		return str
	}
	bytes := strconv.S2B(str)
	l := len(bytes)
	b := bytes[l-1]
	buffer := make([]byte, 0, l)
	buffer = append(buffer, bytes...)
	if b == 's' || b == 'S' {
		buffer = append(buffer, 'e')
	}
	buffer = append(buffer, 's')
	return strconv.B2S(buffer)
}
