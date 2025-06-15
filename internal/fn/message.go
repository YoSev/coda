package fn

import (
	"encoding/json"
	"errors"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type Shoutrrr struct {
	Urls       []string      `json:"urls" yaml:"urls"`
	Message    string        `json:"message" yaml:"message"`
	Parameters *types.Params `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (f *Fn) Shoutrrr(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *Shoutrrr) (json.RawMessage, error) {
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
