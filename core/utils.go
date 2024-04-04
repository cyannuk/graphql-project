package core

import (
	"context"
	"net"
	"net/netip"
	"strconv"
	"strings"

	gotils "github.com/savsgio/gotils/strconv"
)

func TrimQuotes(str string) string {
	if str == "" {
		return str
	}
	bytes := gotils.S2B(str)
	first := bytes[0]
	l := len(bytes) - 1
	last := bytes[l]
	if (first == '\'' && last == '\'') || (first == '"' && last == '"') || (first == '`' && last == '`') {
		return gotils.B2S(bytes[1:l])
	}
	return str
}

func StartWith(str string, ch byte) bool {
	if str == "" {
		return false
	}
	for _, b := range gotils.S2B(str) {
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
	bytes := gotils.S2B(str)
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
	bb := gotils.S2B(str)
	b := bb[0]
	if !(b >= 'a' && b <= 'z') {
		return str
	}
	var i int
	for i, b = range bb {
		if b != '_' {
			break
		}
	}
	if b >= 'a' && b <= 'z' {
		buffer := make([]byte, 0, len(bb))
		buffer = append(buffer, b-0x20)
		buffer = append(buffer, bb[i+1:]...)
		return gotils.B2S(buffer)
	} else {
		if i == 0 {
			return str
		} else {
			return gotils.B2S(bb[i:])
		}
	}
}

// ASCII only
func Uncapitalize(str string) string {
	if str == "" {
		return str
	}
	bb := gotils.S2B(str)
	b := bb[0]
	if !(b >= 'A' && b <= 'Z') {
		return str
	}
	var i int
	for i, b = range bb {
		if b != '_' {
			break
		}
	}
	if b >= 'A' && b <= 'Z' {
		buffer := make([]byte, 0, len(bb))
		buffer = append(buffer, b+0x20)
		buffer = append(buffer, bb[i+1:]...)
		return gotils.B2S(buffer)
	} else {
		if i == 0 {
			return str
		} else {
			return gotils.B2S(bb[i:])
		}
	}
}

// ASCII only
func Plural(str string) string {
	if str == "" {
		return str
	}
	bytes := gotils.S2B(str)
	l := len(bytes)
	b := bytes[l-1]
	buffer := make([]byte, 0, l)
	buffer = append(buffer, bytes...)
	if b == 's' || b == 'S' {
		buffer = append(buffer, 'e')
	}
	buffer = append(buffer, 's')
	return gotils.B2S(buffer)
}

func Quote(s string) string {
	l := len(s)
	if l == 0 {
		return `""`
	}
	b := gotils.S2B(s)
	result := make([]byte, 0, l+2)
	result = append(result, '"')
	result = append(result, b...)
	result = append(result, '"')
	return gotils.B2S(result)
}

func Replace(str string, oldStr string, newStr string, count int) string {
	if oldStr == newStr || count == 0 {
		return str
	}
	return replace(str, oldStr, newStr, &count)
}

func replace(str string, oldStr string, newStr string, count *int) string {
	if *count == 0 {
		return str
	}
	*count--
	i := strings.Index(str, oldStr)
	if i < 0 {
		return str
	}
	return str[:i] + newStr + replace(str[i+len(oldStr):], oldStr, newStr, count)
}

func Join(values ...any) string {
	b := make([]byte, 0, 512)
	for _, value := range values {
		switch v := value.(type) {
		case string:
			b = append(b, v...)
		case int:
			b = strconv.AppendInt(b, int64(v), 10)
		case int8:
			b = strconv.AppendInt(b, int64(v), 10)
		case int16:
			b = strconv.AppendInt(b, int64(v), 10)
		case int32:
			b = strconv.AppendInt(b, int64(v), 10)
		case int64:
			b = strconv.AppendInt(b, v, 10)
		case uint:
			b = strconv.AppendUint(b, uint64(v), 10)
		case uint8:
			b = strconv.AppendUint(b, uint64(v), 10)
		case uint16:
			b = strconv.AppendUint(b, uint64(v), 10)
		case uint32:
			b = strconv.AppendUint(b, uint64(v), 10)
		case uint64:
			b = strconv.AppendUint(b, v, 10)
		case float32:
			b = strconv.AppendFloat(b, float64(v), 'f', -1, 32)
		case float64:
			b = strconv.AppendFloat(b, v, 'f', -1, 64)
		case bool:
			b = strconv.AppendBool(b, v)
		case []byte:
			b = append(b, v...)
		}
	}
	return gotils.B2S(b)
}

func AppendStrings(b []byte, values ...string) []byte {
	for _, v := range values {
		b = append(b, v...)
	}
	return b
}

func ParseAddress(addr string) (ip netip.Addr, err error) {
	ip, err = netip.ParseAddr(addr)
	if err == nil {
		return
	}
	addresses, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip4", addr)
	if err == nil {
		ip = addresses[0]
	}
	return
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
