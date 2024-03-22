package tests

import (
	"errors"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"

	"graphql-project/core"
)

const idMargin = 5
const timestampMargin = time.Second * 5

func Include(fields ...string) func(cmp.Path) bool {
	return func(p cmp.Path) bool {
		if v, ok := p.Last().(cmp.MapIndex); ok {
			return slices.Contains(fields, v.Key().String())
		}
		return false
	}
}

func Exclude(fields ...string) func(cmp.Path) bool {
	return func(p cmp.Path) bool {
		if v, ok := p.Last().(cmp.MapIndex); ok {
			return !slices.Contains(fields, v.Key().String())
		}
		return true
	}
}

func isBasicType(v any) bool {
	switch v.(type) {
	case time.Time, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, string:
		return true
	default:
		return false
	}
}

func BasicTypes(x, y any) bool {
	return isBasicType(x) && isBasicType(y)
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t.UTC(), nil
	} else {
		if t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999", s, time.UTC); err == nil {
			return t, nil
		} else {
			return time.ParseInLocation(time.DateOnly, s, time.UTC)
		}
	}
}

func getTime(v any) (time.Time, error) {
	switch t := v.(type) {
	case time.Time:
		return t, nil
	case string:
		return parseTime(t)
	default:
		return time.Time{}, errors.New("not time")
	}
}

func TimestampCompare(diff time.Duration) func(x, y any) bool {
	return func(x, y any) bool {
		if t1, err := getTime(x); err == nil {
			if t2, err := getTime(y); err == nil {
				return t1.Sub(t2).Abs() <= diff
			}
		}
		return x == y
	}
}

func ValueCompare(x, y any) bool {
	if s1, err := getString(x); err == nil {
		if s2, err := getString(y); err == nil {
			return s1 == s2
		}
	}
	if i1, err := getInt(x); err == nil {
		if i2, err := getInt(y); err == nil {
			return i1 == i2
		}
	}
	if b1, err := getBool(x); err == nil {
		if b2, err := getBool(y); err == nil {
			return b1 == b2
		}
	}
	if f1, err := getFloat(x); err == nil {
		if f2, err := getFloat(y); err == nil {
			return f1 == f2
		}
	}
	if t1, err := getTime(x); err == nil {
		if t2, err := getTime(y); err == nil {
			return t1.Compare(t2) == 0
		}
	}
	return false
}

func getFloat(v any) (float64, error) {
	switch f := v.(type) {
	case int:
		return float64(f), nil
	case int8:
		return float64(f), nil
	case int16:
		return float64(f), nil
	case int32:
		return float64(f), nil
	case int64:
		return float64(f), nil
	case uint:
		return float64(f), nil
	case uint8:
		return float64(f), nil
	case uint16:
		return float64(f), nil
	case uint32:
		return float64(f), nil
	case uint64:
		return float64(f), nil
	case float32:
		return float64(f), nil
	case float64:
		return f, nil
	case string:
		return strconv.ParseFloat(f, 64)
	default:
		return 0, errors.New("not float")
	}
}

func getInt(v any) (int64, error) {
	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int64:
		return i, nil
	case uint:
		return int64(i), nil
	case uint8:
		return int64(i), nil
	case uint16:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint64:
		return int64(i), nil
	case string:
		return strconv.ParseInt(i, 10, 64)
	default:
		return 0, errors.New("not int")
	}
}

func getBool(v any) (bool, error) {
	switch b := v.(type) {
	case bool:
		return b, nil
	case string:
		return strconv.ParseBool(b)
	default:
		return false, errors.New("not bool")
	}
}

func getString(v any) (string, error) {
	switch s := v.(type) {
	case string:
		return s, nil
	default:
		return "", errors.New("not string")
	}
}

func IdCompare(diff int64) func(x, y any) bool {
	return func(x, y any) bool {
		if i1, err := getInt(x); err == nil {
			if i2, err := getInt(y); err == nil {
				return core.Abs(i1-i2) <= diff
			}
		}
		if f1, err := getFloat(x); err == nil {
			if f2, err := getFloat(y); err == nil {
				return math.Abs(f1-f2) <= float64(diff)
			}
		}
		return x == y
	}
}

func PasswordCompare(x, y any) bool {
	if s1, err := getString(x); err == nil {
		if s2, err := getString(y); err == nil {
			if s1 == s2 || core.VerifyPassword(s1, s2) {
				return true
			}
		}
	}
	return false
}

func compare(expected map[string]any, actual map[string]any) error {
	// idCompare := cmp.FilterPath(Include("id"), cmp.Comparer(IdCompare(5)))
	passwordCompare := cmp.FilterPath(Include("password"), cmp.Comparer(PasswordCompare))
	timestampCompare := cmp.FilterPath(Include("createdAt", "deletedAt"), cmp.Comparer(TimestampCompare(timestampMargin)))
	valueCompare := cmp.FilterPath(Exclude("password", "createdAt", "deletedAt"), cmp.FilterValues(BasicTypes, cmp.Comparer(ValueCompare)))
	var r reporter
	if cmp.Equal(expected, actual, passwordCompare, timestampCompare, valueCompare, cmp.Reporter(&r)) {
		return nil
	}
	return errors.New(r.String())
}
