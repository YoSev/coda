package fn

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (f *fnCompare) compare(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *compareParams) (json.RawMessage, error) {
		var result bool
		switch params.Operator {
		case OpEq:
			cmp, err := comparePrimitive(params.Left, params.Right)
			if err != nil {
				return nil, err
			}
			result = cmp == 0

		case OpNe:
			cmp, err := comparePrimitive(params.Left, params.Right)
			if err != nil {
				return nil, err
			}
			result = cmp != 0

		case OpGt, OpGte, OpLt, OpLte:
			cmp, err := comparePrimitive(params.Left, params.Right)
			if err != nil {
				return nil, err
			}

			switch params.Operator {
			case OpGt:
				result = cmp > 0
			case OpGte:
				result = cmp >= 0
			case OpLt:
				result = cmp < 0
			case OpLte:
				result = cmp <= 0
			}

		case OpContains:
			ls := fmt.Sprint(params.Left)
			rs := fmt.Sprint(params.Right)
			result = strings.Contains(ls, rs)

		case OpEmpty:
			result = params.Left == nil || params.Left == ""

		case OpNotEmpty:
			result = !(params.Left == nil || params.Left == "")

		default:
			return nil, fmt.Errorf("unknown operator: %s", params.Operator)
		}

		if !result {
			return nil, errors.New("comparison failed")
		}

		return utils.ReturnRaw(result), nil
	})
}

func comparePrimitive(a, b any) (int, error) {
	// nil
	if a == nil || b == nil {
		switch {
		case a == nil && b == nil:
			return 0, nil
		case a == nil:
			return -1, nil
		default:
			return 1, nil
		}
	}

	switch x := a.(type) {

	case bool:
		y, ok := b.(bool)
		if !ok {
			return 0, fmt.Errorf("type mismatch")
		}
		if x == y {
			return 0, nil
		}
		if !x && y {
			return -1, nil
		}
		return 1, nil

	case string:
		y, ok := b.(string)
		if !ok {
			return 0, fmt.Errorf("type mismatch")
		}
		return strings.Compare(x, y), nil

	case float64:
		y, ok := b.(float64)
		if !ok {
			return 0, fmt.Errorf("type mismatch")
		}
		switch {
		case x < y:
			return -1, nil
		case x > y:
			return 1, nil
		default:
			return 0, nil
		}
	}

	return 0, fmt.Errorf("unsupported type %T", a)
}
