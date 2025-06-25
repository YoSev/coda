package fn

import (
	"encoding/json"
	"strings"

	"github.com/iancoleman/strcase"
)

type FnString struct {
	FnHandler
	*Fn

	Category FnCategory
}

func (f *FnString) Register() {
	f.fns["string.upper"] = &FnEntry{
		Handler:     f.upperCase,
		Name:        "Uppercase",
		Description: "Converts a string to uppercase",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to uppercase", Mandatory: true},
		},
	}

	f.fns["string.lower"] = &FnEntry{
		Handler:     f.dowerCase,
		Name:        "Lowercase",
		Description: "Converts a string to lowercase",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to lowercase", Mandatory: true},
		},
	}

	f.fns["string.camel"] = &FnEntry{
		Handler:     f.camelCase,
		Name:        "Camel Case",
		Description: "Converts a string to camel case",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to camel case", Mandatory: true},
		},
	}

	f.fns["string.snake"] = &FnEntry{
		Handler:     f.snakeCase,
		Name:        "Snake Case",
		Description: "Converts a string to snake case",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to snake case", Mandatory: true},
		},
	}

	f.fns["string.kebap"] = &FnEntry{
		Handler:     f.kebapCase,
		Name:        "Kebab Case",
		Description: "Converts a string to kebab case",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to convert to kebab case", Mandatory: true},
		},
	}

	f.fns["string.reverse"] = &FnEntry{
		Handler:     f.stringReverse,
		Name:        "Reverse String",
		Description: "Reverses the characters in a string",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to reverse", Mandatory: true},
		},
	}

	f.fns["string.trim"] = &FnEntry{
		Handler:     f.stringTrim,
		Name:        "Trim String",
		Description: "Trims whitespace or specified characters from the start and end of a string",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to trim", Mandatory: true},
			{Name: "delimiter", Description: "The characters to trim (default is whitespace)", Mandatory: false},
		},
	}

	f.fns["string.split"] = &FnEntry{
		Handler:     f.stringSplit,
		Name:        "Split String",
		Description: "Splits a string into an array based on a delimiter",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string to split", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to use for splitting (default is whitespace)", Mandatory: false},
		},
	}

	f.fns["string.join"] = &FnEntry{
		Handler:     f.stringJoin,
		Name:        "Join Strings",
		Description: "Joins an array of strings into a single string using a delimiter",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The array of strings to join", Type: "array", Mandatory: true},
			{Name: "delimiter", Description: "The delimiter to use for joining (default is ',')", Mandatory: false},
		},
	}

	f.fns["string.resolve"] = &FnEntry{
		Handler:     f.stringResolve,
		Name:        "Resolve String",
		Description: "Resolves a string value, useful for dynamic values",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to resolve", Type: "any", Mandatory: true},
		},
	}

	f.fns["string"] = &FnEntry{
		Handler:     f.string,
		Name:        "String",
		Description: "Returns the string value as is",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The string value to return", Type: "string", Mandatory: true},
		},
	}

	f.fns["json.encode"] = &FnEntry{
		Handler:     f.jsonEncode,
		Name:        "JSON Encode",
		Description: "Encodes a value to a JSON string",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to encode as JSON", Type: "any", Mandatory: true},
		},
	}

	f.fns["json.decode"] = &FnEntry{
		Handler:     f.jsonDecode,
		Name:        "JSON Decode",
		Description: "Decodes a JSON string into a value",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The JSON string to decode", Type: "string", Mandatory: true},
		},
	}
}

type anyParams struct {
	Value any `json:"value" yaml:"value"`
}

type stringParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *FnString) upperCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return json.RawMessage(addQuotes(strings.ToUpper(params.Value))), nil
	})
}

func (f *FnString) dowerCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strings.ToLower(params.Value)), nil
	})
}

func (f *FnString) camelCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strcase.ToCamel(params.Value)), nil
	})
}

func (f *FnString) snakeCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strcase.ToSnake(params.Value)), nil
	})
}

func (f *FnString) kebapCase(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(strcase.ToKebab(params.Value)), nil
	})
}

func (f *FnString) stringReverse(j json.RawMessage) (json.RawMessage, error) {
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

func (f *FnString) stringTrim(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		delimiter := params.Delimiter
		if params.Delimiter == "" {
			delimiter = " "
		}
		return returnRaw(strings.Trim(params.Value, delimiter)), nil
	})
}

func (f *FnString) stringSplit(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringWithDelimiterParams) (json.RawMessage, error) {
		return returnRaw(strings.Split(params.Value, params.Delimiter)), nil
	})
}

type stringArrayWithDelimiterParams struct {
	Value     []string `json:"value" yaml:"value"`
	Delimiter string   `json:"delimiter" yaml:"delimiter"`
}

func (f *FnString) stringJoin(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringArrayWithDelimiterParams) (json.RawMessage, error) {
		return returnRaw(strings.Join(params.Value, params.Delimiter)), nil
	})
}

func (f *FnString) stringResolve(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *anyParams) (json.RawMessage, error) {
		return returnRaw(params.Value), nil
	})
}

func (f *FnString) string(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		return returnRaw(params.Value), nil
	})
}

func (f *FnString) jsonEncode(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *anyParams) (json.RawMessage, error) {
		out, err := json.Marshal(params.Value)
		if err != nil {
			return nil, err
		}
		return returnRaw(out), nil
	})
}

func (f *FnString) jsonDecode(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stringParams) (json.RawMessage, error) {
		var out any
		if err := json.Unmarshal([]byte(params.Value), &out); err != nil {
			return nil, err
		}
		return returnRaw(out), nil
	})
}
