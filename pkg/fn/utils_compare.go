package fn

import (
	"encoding/json"
	"errors"
	"fmt"
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
		Description: "Compares two values using various operators (eq, gt, lt, contains, empty).",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "left", Mandatory: true},
			{Name: "operator", Mandatory: true},
			{Name: "right", Mandatory: false},
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
	OpEmpty    CompareOperator = "empty"
	OpNotEmpty CompareOperator = "not_empty"
)

type compareParams struct {
	Left     any             `json:"left"`
	Operator CompareOperator `json:"operator"`
	Right    any             `json:"right,omitempty"`
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
			result = equal(left, right)

		case OpNe:
			result = !equal(left, right)

		case OpGt:
			result = compareNumbers(left, right, func(a, b float64) bool { return a > b })

		case OpGte:
			result = compareNumbers(left, right, func(a, b float64) bool { return a >= b })

		case OpLt:
			result = compareNumbers(left, right, func(a, b float64) bool { return a < b })

		case OpLte:
			result = compareNumbers(left, right, func(a, b float64) bool { return a <= b })

		case OpContains:
			result = strings.Contains(fmt.Sprintf("%v", left), fmt.Sprintf("%v", right))

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

func equal(a, b any) bool {
	// fast path
	switch x := a.(type) {

	case float64:
		y, ok := toNumber(b)
		return ok && x == y

	case string:
		y, ok := b.(string)
		return ok && x == y

	case bool:
		y, ok := b.(bool)
		return ok && x == y
	}

	// fallback
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func compareNumbers(a any, b any, cmp func(float64, float64) bool) bool {
	x, ok := toNumber(a)
	if !ok {
		return false
	}

	y, ok := toNumber(b)
	if !ok {
		return false
	}

	return cmp(x, y)
}

func toNumber(v any) (float64, bool) {
	switch n := v.(type) {

	case float64:
		return n, true

	case float32:
		return float64(n), true

	case int:
		return float64(n), true

	case int64:
		return float64(n), true

	case int32:
		return float64(n), true

	case uint:
		return float64(n), true

	case json.Number:
		f, err := n.Float64()
		return f, err == nil

	case string:
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil

	default:
		return 0, false
	}
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}

	switch x := v.(type) {

	case string:
		return strings.TrimSpace(x) == ""

	case []any:
		return len(x) == 0

	case map[string]any:
		return len(x) == 0

	case json.RawMessage:
		return len(x) == 0

	case bool:
		return !x

	case float64:
		return x == 0

	case int:
		return x == 0

	default:
		return fmt.Sprintf("%v", v) == ""
	}
}
