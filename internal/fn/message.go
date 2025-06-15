package fn

import (
	"encoding/json"
	"errors"

	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type Shoutrrr struct {
	Url        []string      `json:"url,omitempty" yaml:"url,omitempty"`
	Message    string        `json:"message,omitempty" yaml:"message,omitempty"`
	Parameters *types.Params `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (f *Fn) Shoutrrr(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *Shoutrrr) (json.RawMessage, error) {
		sender, err := shoutrrr.CreateSender(params.Url...)
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
