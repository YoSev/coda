package fn

import (
	"encoding/json"
	"math"
)

type mathStruct struct {
	Value  float64 `json:"value" yaml:"value"`
	Amount float64 `json:"amount" yaml:"amount"`
}

func (f *Fn) MathInc(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value += params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *Fn) MathDec(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value -= params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *Fn) MathMultiply(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value *= params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *Fn) MathDivide(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		params.Amount = math.Max(1, params.Amount)
		params.Value /= params.Amount
		return returnRaw(params.Value), nil
	})
}

func (f *Fn) MathModulo(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *mathStruct) (json.RawMessage, error) {
		return returnRaw(math.Mod(params.Value, params.Amount)), nil
	})
}
