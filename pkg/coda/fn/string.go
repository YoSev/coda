package fn

import (
	"encoding/json"
	"strings"

	"github.com/iancoleman/strcase"
)

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
		return returnAny(strings.ToLower(params.Value)), nil
	})
}

func (f *Fn) CamelCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnAny(strcase.ToCamel(params.Value)), nil
	})
}

func (f *Fn) SnakeCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnAny(strcase.ToSnake(params.Value)), nil
	})
}

func (f *Fn) KebapCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnAny(strcase.ToKebab(params.Value)), nil
	})
}

func (f *Fn) StringReverse(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		runes := []rune(params.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return returnAny(string(runes)), nil
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
		return returnAny(strings.Trim(params.Value, delimiter)), nil
	})
}

func (f *Fn) StringSplit(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		return returnAny(strings.Split(params.Value, params.Delimiter)), nil
	})
}

type stringArrayWithDelimiterParams struct {
	Value     []string `json:"value" yaml:"value"`
	Delimiter string   `json:"delimiter" yaml:"delimiter"`
}

func (f *Fn) StringJoin(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringArrayWithDelimiterParams) (json.RawMessage, error) {
		return returnAny(strings.Join(params.Value, params.Delimiter)), nil
	})
}
