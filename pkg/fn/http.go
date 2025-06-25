package fn

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/yosev/coda/internal/utils"
)

type fnHttp struct {
	*Fn
	category FnCategory
}

func (f *fnHttp) init(fn *Fn) {
	f.Fn = fn

	fn.register("http.request", &FnEntry{
		Handler:     f.httpReq,
		Name:        "HTTP Request",
		Description: "Performs an HTTP request",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "url", Description: "The url to query", Mandatory: true},
			{Name: "method", Description: "The HTTP method to use", Enum: []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, Mandatory: true},
			{Name: "headers", Description: "The Headers to use", Type: "object", Mandatory: false},
			{Name: "body", Description: "The Body to use", Type: "any", Mandatory: false},
		},
	})
}

type HttpReqParams struct {
	Url     string            `json:"url" yaml:"url"`
	Method  string            `json:"method" yaml:"method"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    any               `json:"body" yaml:"body"`
}

func (f *fnHttp) httpReq(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *HttpReqParams) (json.RawMessage, error) {
		client := resty.New()

		request := client.R()
		request.SetBody(params.Body)
		request.SetHeaders(params.Headers)

		var response *resty.Response
		var err error

		if _, ok := params.Headers["User-Agent"]; !ok {
			request.SetHeader("User-Agent", fmt.Sprintf("coda/%s", f.Fn.version))
		}

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

		return utils.ReturnRaw(resp), nil
	})
}
