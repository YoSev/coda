package fn

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/yosev/coda/internal/utils"
)

type fnCompare struct {
	category FnCategory
}

func (f *fnCompare) init(fn *Fn) {
	fn.register("utils.compare", &FnEntry{
		Handler:     f.compare,
		Name:        "Compare values with various operators",
		Description: "Compares two values using various operators (e.g. eq, gt, contains) and returns true if the comparison holds.",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "left", Description: "The left value to compare", Mandatory: true},
			{Name: "operator", Description: "The comparison operator", Mandatory: true},
			{Name: "right", Description: "The right value to compare", Mandatory: false},
		},
	})
}

type CompareOperator string

const (
	OpEq CompareOperator = "eq"
	OpNe CompareOperator = "ne"

	OpGt  CompareOperator = "gt"
	OpGte CompareOperator = "gte"
	OpLt  CompareOperator = "lt"
	OpLte CompareOperator = "lte"

	OpContains CompareOperator = "contains"

	OpIn    CompareOperator = "in"
	OpNotIn CompareOperator = "not_in"

	OpEmpty    CompareOperator = "empty"
	OpNotEmpty CompareOperator = "not_empty"
)

type compareParams struct {
	Left     any             `json:"left" yaml:"left"`
	Operator CompareOperator `json:"operator" yaml:"operator"`
	Right    any             `json:"right,omitempty" yaml:"right,omitempty"`
}

func normalize(v any) any {
	switch t := v.(type) {
	case json.Number:
		if f, err := t.Float64(); err == nil {
			return f
		}
	}
	return v
}

func (f *fnCompare) compare(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *compareParams) (json.RawMessage, error) {

		left := normalize(params.Left)
		right := normalize(params.Right)

		var result bool

		switch params.Operator {

		case OpEq:
			result = reflect.DeepEqual(left, right)

		case OpNe:
			result = !reflect.DeepEqual(left, right)

		case OpGt:
			result = compareNumbers(left, right, func(a, b float64) bool {
				return a > b
			})

		case OpGte:
			result = compareNumbers(left, right, func(a, b float64) bool {
				return a >= b
			})

		case OpLt:
			result = compareNumbers(left, right, func(a, b float64) bool {
				return a < b
			})

		case OpLte:
			result = compareNumbers(left, right, func(a, b float64) bool {
				return a <= b
			})

		case OpContains:
			result = strings.Contains(
				fmt.Sprintf("%v", left),
				fmt.Sprintf("%v", right),
			)

		case OpEmpty:
			result = isEmpty(left)

		case OpNotEmpty:
			result = !isEmpty(left)

		default:
			return nil, errors.New("unknown operator")
		}

		if !result {
			return nil, errors.New("comparison failed")
		}

		return utils.ReturnRaw(true), nil
	})
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}

	switch x := v.(type) {

	case string:
		return len(strings.TrimSpace(x)) == 0

	case []byte:
		return len(x) == 0

	case []any:
		return len(x) == 0

	case map[string]any:
		return len(x) == 0

	case map[any]any:
		return len(x) == 0

	case json.RawMessage:
		return len(x) == 0

	case *string:
		return x == nil || len(strings.TrimSpace(*x)) == 0

	case *int, *int64, *float64, *bool:
		// pointers to scalars are considered empty if nil
		return true

	case int:
		return x == 0

	case int64:
		return x == 0

	case float64:
		return x == 0

	case bool:
		return x == false

	default:
		// fallback: reflect-based emptiness check
		rv := reflect.ValueOf(v)

		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			return rv.IsNil()

		case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
			return rv.Len() == 0

		case reflect.String:
			return strings.TrimSpace(rv.String()) == ""

		case reflect.Bool:
			return !rv.Bool()

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return rv.Int() == 0

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return rv.Uint() == 0

		case reflect.Float32, reflect.Float64:
			return rv.Float() == 0

		default:
			// unknown types: treat zero-value structs as empty
			return reflect.DeepEqual(v, reflect.Zero(rv.Type()).Interface())
		}
	}
}

func compareNumbers(left any, right any, cmp func(a, b float64) bool) bool {
	l, ok := toFloat64(left)
	if !ok {
		return false
	}

	r, ok := toFloat64(right)
	if !ok {
		return false
	}

	return cmp(l, r)
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {

	case float64:
		return n, true

	case float32:
		return float64(n), true

	case int:
		return float64(n), true

	case int8:
		return float64(n), true

	case int16:
		return float64(n), true

	case int32:
		return float64(n), true

	case int64:
		return float64(n), true

	case uint:
		return float64(n), true

	case uint8:
		return float64(n), true

	case uint16:
		return float64(n), true

	case uint32:
		return float64(n), true

	case uint64:
		// beware: precision loss for large uint64, but typical for JSON comparisons
		return float64(n), true

	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return f, true

	case string:
		// optional but often useful if JSON comes in as strings
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return 0, false
		}
		return f, true

	default:
		return 0, false
	}
}
