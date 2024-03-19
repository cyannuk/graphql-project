package tests

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"slices"
	"strconv"
	"text/template"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"
	gotils "github.com/savsgio/gotils/strconv"
	"gopkg.in/yaml.v3"

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

func getTime(v any) (time.Time, error) {
	switch t := v.(type) {
	case time.Time:
		return t, nil
	case string:
		return time.Parse(time.RFC3339Nano, t)
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
	if t1, err := getTime(x); err == nil {
		if t2, err := getTime(y); err == nil {
			return t1.Compare(t2) == 0
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
	if s1, err := getString(x); err == nil {
		if s2, err := getString(y); err == nil {
			return s1 == s2
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

func Now() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func getTemplate(name string) ([]byte, error) {
	if _, err := os.Stat(name); err == nil {
		return parseTemplate(name)
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else {
		return nil, err
	}
}

func parseTemplate(fileName string) ([]byte, error) {
	funcMap := template.FuncMap{
		"NOW": Now,
	}
	tmpl, err := template.New(path.Base(fileName)).Funcs(funcMap).ParseFiles(fileName)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.Grow(1024)
	if err := tmpl.Execute(&buffer, nil); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func loadTestData(name string) (m map[string]any, err error) {
	var data []byte
	if data, err = getTemplate(name + ".yml"); err != nil {
		return
	}
	if data == nil {
		if data, err = getTemplate(name + ".yaml"); err != nil {
			return
		}
		if data == nil {
			if data, err = getTemplate(name + ".json"); err != nil {
				return
			}
			if data == nil {
				err = fmt.Errorf("testdata not found `%s.[yml|yaml|json]`", name)
			} else {
				err = json.Unmarshal(data, &m)
			}
		} else {
			err = yaml.Unmarshal(data, &m)
		}
	} else {
		err = yaml.Unmarshal(data, &m)
	}
	return
}

type gqlQuery struct {
	Query string `json:"query"`
}

func loadRequestData(name string) ([]byte, error) {
	query, err := parseTemplate(name + ".gql")
	if err != nil {
		return nil, err
	}
	return json.Marshal(gqlQuery{gotils.B2S(query)})
}

func Compare(expectedDataFile string, actual map[string]any) error {
	expected, err := loadTestData(expectedDataFile)
	if err != nil {
		return err
	}
	// idCompare := cmp.FilterPath(Include("id"), cmp.Comparer(IdCompare(5)))
	timestampCompare := cmp.FilterPath(Include("createdAt", "deletedAt"), cmp.Comparer(TimestampCompare(timestampMargin)))
	valueCompare := cmp.FilterPath(Exclude("createdAt", "deletedAt"), cmp.FilterValues(BasicTypes, cmp.Comparer(ValueCompare)))
	var r reporter
	if cmp.Equal(expected, actual, timestampCompare, valueCompare, cmp.Reporter(&r)) {
		return nil
	}
	return errors.New(r.String())
}
