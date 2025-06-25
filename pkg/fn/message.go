package fn

import (
	"encoding/json"
	"errors"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type FnMessage struct {
	FnHandler
	*Fn

	Category FnCategory
}

func (f *FnMessage) Register() {
	f.fns["message.shoutrrr"] = &FnEntry{
		Handler:     f.shoutrrr,
		Name:        "Shoutrrr",
		Description: "Sends a message using the Shoutrrr notification system",
		Category:    f.Category,
		Parameters: []FnParameter{
			{Name: "urls", Description: "The shoutrrr targets", Type: "array", Mandatory: true},
			{Name: "message", Description: "The shoutrrr message to send", Mandatory: true},
			{Name: "parameters", Description: "Additional shoutrrr properties", Type: "object", Mandatory: false},
		},
	}
}

type shoutrrrStruct struct {
	Urls       []string      `json:"urls" yaml:"urls"`
	Message    string        `json:"message" yaml:"message"`
	Parameters *types.Params `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (f *FnMessage) shoutrrr(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *shoutrrrStruct) (json.RawMessage, error) {
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
