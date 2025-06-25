package fn

import (
	"encoding/json"
	"math"
)

type FnMath struct {
	FnHandler
	*Fn

	Category FnCategory
}

func (f *FnMath) Register() {
	f.fns["math.inc"] = &FnEntry{
		Handler:     f.inc,
		Name:        "Increment",
		Description: "Increment a value by a specified amount",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to increment", Mandatory: true, Type: "number"},
			{Name: "amount", Description: "The amount to increment by (default to 1)", Mandatory: false, Type: "number"},
		},
	}

	f.fns["math.dec"] = &FnEntry{
		Handler:     f.dec,
		Name:        "Decrement",
		Description: "Decrement a value by a specified amount",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to decrement", Mandatory: true, Type: "number"},
			{Name: "amount", Description: "The amount to decrement by (default to 1)", Mandatory: false, Type: "number"},
		},
	}

	f.fns["math.multiply"] = &FnEntry{
		Handler:     f.multiply,
		Name:        "Multiply",
		Description: "Multiply a value by a specified amount",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to multiply", Mandatory: true, Type: "number"},
			{Name: "amount", Description: "The amount to multiply by (default to 1)", Mandatory: false, Type: "number"},
		},
	}

	f.fns["math.divide"] = &FnEntry{
		Handler:     f.divide,
		Name:        "Divide",
		Description: "Divide a value by a specified amount",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to divide", Mandatory: true, Type: "number"},
			{Name: "amount", Description: "The amount to divide by (default to 1)", Mandatory: false, Type: "number"},
		},
	}

	f.fns["math.modulo"] = &FnEntry{
		Handler:     f.modulo,
		Name:        "Modulo",
		Description: "Calculate the modulo of a value with a specified amount",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The source value", Type: "number", Mandatory: true},
			{Name: "amount", Description: "The amount by which to mod", Type: "number", Mandatory: false},
		},
	}
}

type mathStruct struct {
	Value  float64 `json:"value" yaml:"value"`
	Amount float64 `json:"amount" yaml:"amount"`
}

func (f *FnMath) inc(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value += params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *FnMath) dec(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value -= params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *FnMath) multiply(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value *= params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *FnMath) divide(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value /= params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *FnMath) modulo(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		return returnRaw(math.Mod(params.Value, params.Amount)), nil
	})
}
