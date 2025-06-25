package fn

import (
	"encoding/json"
	"fmt"
	"os"
)

type FnIo struct {
	FnHandler
	*Fn

	Category FnCategory
}

func (f *FnIo) Register() {
	f.fns["io.stdout"] = &FnEntry{
		Handler:     f.stdout,
		Name:        "Write to stdout",
		Description: "Writes a string to stdout",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to write to stdout", Mandatory: true},
		},
	}

	f.fns["io.stderr"] = &FnEntry{
		Handler:     f.stderr,
		Name:        "Write to stderr",
		Description: "Writes a string to stderr",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to write to stderr", Mandatory: true, Type: "string"},
		},
	}
}

type stdParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *FnIo) stdout(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stdout, params.Value)
		return nil, nil
	})
}

func (f *FnIo) stderr(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stderr, params.Value)
		return nil, nil
	})
}
