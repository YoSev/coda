package fn

import (
	"encoding/json"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/yosev/coda/internal/utils"
)

type fnString struct {
	category FnCategory
}

func (f *fnString) init(fn *Fn) {
	fn.register("string.upper", &FnEntry{
		Handler:     f.upperCase,
		Name:        "Uppercase",
		Description: "Converts a string to uppercase",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to uppercase", Mandatory: true},
		},
	})

	fn.register("string.lower", &FnEntry{
		Handler:     f.lowerCase,
		Name:        "Lowercase",
		Description: "Converts a string to lowercase",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to lowercase", Mandatory: true},
		},
	})

	fn.register("string.camel", &FnEntry{
		Handler:     f.camelCase,
		Name:        "Camel Case",
		Description: "Converts a string to camel case",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to camel case", Mandatory: true},
		},
	})

	fn.register("string.snake", &FnEntry{
		Handler:     f.snakeCase,
		Name:        "Snake Case",
		Description: "Converts a string to snake case",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to snake case", Mandatory: true},
		},
	})

	fn.register("string.kebab", &FnEntry{
		Handler:     f.kebabCase,
		Name:        "Kebab Case",
		Description: "Converts a string to kebab case",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to kebab case", Mandatory: true},
		},
	})

	fn.register("string.reverse", &FnEntry{
		Handler:     f.stringReverse,
		Name:        "Reverse String",
		Description: "Reverses the characters in a string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to reverse", Mandatory: true},
		},
	})

	fn.register("string.trim", &FnEntry{
		Handler:     f.stringTrim,
		Name:        "Trim String",
		Description: "Trims whitespace or specified characters from the start and end of a string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to trim", Mandatory: true},
			{Name: "delimiter", Description: "The characters to trim (default is whitespace)", Mandatory: false},
		},
	})

	fn.register("string.split", &FnEntry{
		Handler:     f.stringSplit,
		Name:        "Split String",
		Description: "Splits a string into an array based on a delimiter",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to split", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to use for splitting (default is whitespace)", Mandatory: false},
		},
	})

	fn.register("string.join", &FnEntry{
		Handler:     f.stringJoin,
		Name:        "Join Strings",
		Description: "Joins an array of strings into a single string using a delimiter",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The array of strings to join", Type: "array", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to use for joining (default is ',')", Mandatory: false},
		},
	})

	fn.register("string.resolve", &FnEntry{
		Handler:     f.stringResolve,
		Name:        "Resolve String",
		Description: "Resolves a string value, useful for dynamic values",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to resolve", Type: "any", Mandatory: true},
		},
	})

	fn.register("string", &FnEntry{
		Handler:     f.string,
		Name:        "String",
		Description: "Returns the string value as is",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string value to return", Type: "string", Mandatory: true},
		},
	})

	fn.register("json.encode", &FnEntry{
		Handler:     f.jsonEncode,
		Name:        "JSON Encode",
		Description: "Encodes a value to a JSON string",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to encode as JSON", Type: "any", Mandatory: true},
		},
	})

	fn.register("json.decode", &FnEntry{
		Handler:     f.jsonDecode,
		Name:        "JSON Decode",
		Description: "Decodes a JSON string into a value",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The JSON string to decode", Type: "string", Mandatory: true},
		},
	})
}

type anyParams struct {
	Value any `json:"value" yaml:"value"`
}

type stringParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *fnString) upperCase(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strings.ToUpper(params.Value)), nil
	})
}

func (f *fnString) lowerCase(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strings.ToLower(params.Value)), nil
	})
}

func (f *fnString) camelCase(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strcase.ToCamel(params.Value)), nil
	})
}

func (f *fnString) snakeCase(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strcase.ToSnake(params.Value)), nil
	})
}

func (f *fnString) kebabCase(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strcase.ToKebab(params.Value)), nil
	})
}

func (f *fnString) stringReverse(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		runes := []rune(params.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return utils.ReturnRaw(string(runes)), nil
	})
}

type stringWithDelimiterParams struct {
	Value     string `json:"value" yaml:"value"`
	Delimiter string `json:"delimiter" yaml:"delimiter"`
}

func (f *fnString) stringTrim(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		delimiter := params.Delimiter
		if params.Delimiter == "" {
			delimiter = " "
		}
		return utils.ReturnRaw(strings.Trim(params.Value, delimiter)), nil
	})
}

func (f *fnString) stringSplit(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strings.Split(params.Value, params.Delimiter)), nil
	})
}

type stringArrayWithDelimiterParams struct {
	Value     []string `json:"value" yaml:"value"`
	Delimiter string   `json:"delimiter" yaml:"delimiter"`
}

func (f *fnString) stringJoin(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringArrayWithDelimiterParams) (json.RawMessage, error) {
		return utils.ReturnRaw(strings.Join(params.Value, params.Delimiter)), nil
	})
}

func (f *fnString) stringResolve(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *anyParams) (json.RawMessage, error) {
		return utils.ReturnRaw(params.Value), nil
	})
}

func (f *fnString) string(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return utils.ReturnRaw(params.Value), nil
	})
}

func (f *fnString) jsonEncode(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *anyParams) (json.RawMessage, error) {
		out, err := json.Marshal(params.Value)
		if err != nil {
			return nil, err
		}
		return utils.ReturnRaw(out), nil
	})
}

func (f *fnString) jsonDecode(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		var out any
		if err := json.Unmarshal([]byte(params.Value), &out); err != nil {
			return nil, err
		}
		return utils.ReturnRaw(out), nil
	})
}
