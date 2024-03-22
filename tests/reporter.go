package tests

import (
	"bufio"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	gotils "github.com/savsgio/gotils/strconv"
	"gopkg.in/yaml.v3"
)

type node struct {
	parent *node
	child  []*node
	name   string
	x      reflect.Value
	y      reflect.Value
}

type reporter struct {
	path    cmp.Path
	root    *node
	current *node
}

func applyIndent(s, indent, mark string) (r string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		r += indent + mark + scanner.Text() + "\n"
	}
	return
}

func toString(value reflect.Value, indent string, mark string) string {
	var s string
	switch value.Kind() {
	case reflect.Invalid:
		return ""
	case reflect.Bool:
		s = strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		s = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.String:
		s = value.String()
	case reflect.Map, reflect.Struct, reflect.Interface, reflect.Slice, reflect.Array:
		if !value.IsNil() {
			b, _ := yaml.Marshal(value.Interface())
			return applyIndent(gotils.B2S(b), indent, mark)
		}
		s = "null"
	default:
		s = "<unsupported value type>"
	}
	return indent + mark + s + "\n"
}

func (r *reporter) PushStep(step cmp.PathStep) {
	if v, ok := step.(cmp.MapIndex); ok {
		n := node{name: v.Key().String()}
		if r.root == nil {
			r.root = &node{name: ""}
			n.parent = r.root
			r.root.child = append(r.root.child, &n)
		} else {
			n.parent = r.current
			r.current.child = append(r.current.child, &n)
		}
		r.current = &n
	}
	r.path = append(r.path, step)
}

func (r *reporter) PopStep() {
	if v, ok := r.path.Last().(cmp.MapIndex); ok {
		if r.current.name == v.Key().String() {
			r.current = r.current.parent
		}
	}
	r.path = r.path[:len(r.path)-1]
}

func (r *reporter) Report(result cmp.Result) {
	if !result.Equal() {
		last := r.path.Last()
		if v, ok := last.(cmp.MapIndex); ok {
			if r.current.name == v.Key().String() {
				r.current.x, r.current.y = last.Values()
			}
		}
	}
}

func indent(i int) string {
	if i > 0 {
		return strings.Repeat("  ", i)
	}
	return ""
}

func addNode(n *node, s string, i int) (string, int) {
	p := indent(i)
	if len(n.child) > 0 {
		if n.name != "" {
			s += p + n.name + "\n"
		}
		i++
		for _, c := range n.child {
			s, i = addNode(c, s, i)
		}
	} else if n.x.Kind() != reflect.Invalid || n.y.Kind() != reflect.Invalid {
		s += p + n.name + ":\n"
		s += toString(n.x, p, "  < ") + toString(n.y, p, "  > ")
	}
	return s, i
}

func (r *reporter) String() string {
	s, _ := addNode(r.root, "", -1)
	return s
}
