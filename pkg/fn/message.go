package fn

import (
	"encoding/json"
	"errors"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/yosev/coda/internal/utils"
)

type fnMessage struct {
	category FnCategory
}

func (f *fnMessage) init(fn *Fn) {
	fn.register("message.shoutrrr", &FnEntry{
		Handler:     f.shoutrrr,
		Name:        "Shoutrrr",
		Description: "Sends a message using the Shoutrrr notification system",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "urls", Description: "The shoutrrr targets", Type: "array", Mandatory: true},
			{Name: "message", Description: "The shoutrrr message to send", Mandatory: true},
			{Name: "parameters", Description: "Additional shoutrrr properties", Type: "object", Mandatory: false},
		},
	})
}

type shoutrrrStruct struct {
	Urls       []string      `json:"urls" yaml:"urls"`
	Message    string        `json:"message" yaml:"message"`
	Parameters *types.Params `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (f *fnMessage) shoutrrr(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *shoutrrrStruct) (json.RawMessage, error) {
		sender, err := shoutrrr.CreateSender(params.Urls...)
		if err != nil {
			return nil, err
		}
		if params.Message == "" {
			return nil, errors.New("message cannot be empty")
		}
		erro := sender.Send(params.Message, params.Parameters)
		if len(erro) > 0 {
			return nil, erro[0]
		}
		return nil, nil
	})
}
