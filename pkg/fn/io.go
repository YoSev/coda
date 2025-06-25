package fn

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yosev/coda/internal/utils"
)

type fnIo struct {
	category FnCategory
}

func (f *fnIo) init(fn *Fn) {
	fn.register("io.stdout", &FnEntry{
		Handler:     f.stdout,
		Name:        "Write to stdout",
		Description: "Writes a string to stdout",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to write to stdout", Mandatory: true},
		},
	})

	fn.register("io.stderr", &FnEntry{
		Handler:     f.stderr,
		Name:        "Write to stderr",
		Description: "Writes a string to stderr",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "value", Description: "The value to write to stderr", Mandatory: true, Type: "string"},
		},
	})
}

type stdParams struct {
	Value string `json:"value" yaml:"value"`
}

func (f *fnIo) stdout(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stdout, params.Value)
		return nil, nil
	})
}

func (f *fnIo) stderr(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *stdParams) (json.RawMessage, error) {
		fmt.Fprintln(os.Stderr, params.Value)
		return nil, nil
	})
}
