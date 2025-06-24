package fn

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type HttpReqParams struct {
	Url     string            `json:"url" yaml:"url"`
	Method  string            `json:"method" yaml:"method"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    any               `json:"body" yaml:"body"`
}

func (f *Fn) HttpReq(j json.RawMessage) (json.RawMessage, error) {
	return handleJSON(j, func(params *HttpReqParams) (json.RawMessage, error) {
		client := resty.New()

		request := client.R()
		request.SetBody(params.Body)
		request.SetHeaders(params.Headers)

		var response *resty.Response
		var err error

		switch strings.ToUpper(params.Method) {
		case "GET":
			response, err = request.Get(params.Url)
		case "POST":
			response, err = request.Post(params.Url)
		case "PUT":
			response, err = request.Put(params.Url)
		case "PATCH":
			response, err = request.Patch(params.Url)
		case "DELETE":
			response, err = request.Delete(params.Url)
		case "HEAD":
			response, err = request.Head(params.Url)
		case "OPTIONS":
			response, err = request.Options(params.Url)
		default:
			return nil, fmt.Errorf("unsupported HTTP method: %s", params.Method)
		}

		if err != nil {
			return nil, fmt.Errorf("error making HTTP request: %w", err)
		}

		resp := map[string]any{
			"status":  response.StatusCode(),
			"headers": response.Header(),
			"body":    string(response.Body()),
		}

		if response.Header().Get("Content-Type") == "application/json" {
			var j map[string]interface{}
			if err := json.Unmarshal(response.Body(), &j); err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
			}
			resp["body"] = j
		}

		if response.StatusCode() >= 400 {
			return nil, fmt.Errorf("HTTP request failed with status %d: %s", response.StatusCode(), string(response.Body()))
		}

		return returnRaw(resp), nil
	})
}
