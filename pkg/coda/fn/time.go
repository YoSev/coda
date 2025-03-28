package fn

import (
	"encoding/json"
	"time"
)

type generateDatetimeParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *Fn) GenerateDatetime(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *generateDatetimeParams) (json.RawMessage, error) {
		t := time.Now().Format(params.Value)
		return returnRaw(t), nil
	})
}
func (f *Fn) GenerateTimestampSec(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().Unix()), nil
}
func (f *Fn) GenerateTimestampMilli(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().UnixMilli()), nil
}
func (f *Fn) GenerateTimestampMicro(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().UnixMicro()), nil
}
func (f *Fn) GenerateTimestampNano(j json.RawMessage) (json.RawMessage, error) {
	return returnRaw(time.Now().UnixNano()), nil
}

type sleepParams struct {
	Value int64 `json:"value" yaml:"value"`
}

func (f *Fn) Sleep(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *sleepParams) (json.RawMessage, error) {
		time.Sleep(time.Duration(params.Value) * time.Millisecond)
		return nil, nil
	})
}
