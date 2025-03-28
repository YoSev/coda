package fn

import (
	"encoding/json"
	"fmt"
	"os"
)

type stdParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *Fn) GetEnv(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		return returnAny(os.Getenv(params.Value)), nil
	})
}

func (f *Fn) Stdout(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stdout, params.Value)
		return nil, nil
	})
}

func (f *Fn) Stderr(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stderr, params.Value)
		return nil, nil
	})
}

type exitParams struct {
	Value int `json:"value" yaml:"value"`
}

func (f *Fn) Exit(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *exitParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stdout, params.Value)
		return nil, nil
	})
}
