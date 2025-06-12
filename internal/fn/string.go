package fn

import (
	"encoding/json"
	"strings"

	"github.com/iancoleman/strcase"
)

type anyParams struct {
	Value any `json:"value" yaml:"value"`
}

type stringParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *Fn) UpperCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return json.RawMessage(addQuotes(strings.ToUpper(params.Value))), nil
	})
}

func (f *Fn) LowerCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strings.ToLower(params.Value)), nil
	})
}

func (f *Fn) CamelCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strcase.ToCamel(params.Value)), nil
	})
}

func (f *Fn) SnakeCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strcase.ToSnake(params.Value)), nil
	})
}

func (f *Fn) KebapCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strcase.ToKebab(params.Value)), nil
	})
}

func (f *Fn) StringReverse(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		runes := []rune(params.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return returnRaw(string(runes)), nil
	})
}

type stringWithDelimiterParams struct {
	Value     string `json:"value" yaml:"value"`
	Delimiter string `json:"delimiter" yaml:"delimiter"`
}

func (f *Fn) StringTrim(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		delimiter := params.Delimiter
		if params.Delimiter == "" {
			delimiter = " "
		}
		return returnRaw(strings.Trim(params.Value, delimiter)), nil
	})
}

func (f *Fn) StringSplit(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		return returnRaw(strings.Split(params.Value, params.Delimiter)), nil
	})
}

type stringArrayWithDelimiterParams struct {
	Value     []string `json:"value" yaml:"value"`
	Delimiter string   `json:"delimiter" yaml:"delimiter"`
}

func (f *Fn) StringJoin(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringArrayWithDelimiterParams) (json.RawMessage, error) {
		return returnRaw(strings.Join(params.Value, params.Delimiter)), nil
	})
}

func (f *Fn) StringResolve(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *anyParams) (json.RawMessage, error) {
		return returnRaw(params.Value), nil
	})
}

func (f *Fn) JsonEncode(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *anyParams) (json.RawMessage, error) {
		out, err := json.Marshal(params.Value)
		if err != nil {
			return nil, err
		}
		return returnRaw(out), nil
	})
}

func (f *Fn) JsonDecode(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		var out any
		if err := json.Unmarshal([]byte(params.Value), &out); err != nil {
			return nil, err
		}
		return returnRaw(out), nil
	})
}
