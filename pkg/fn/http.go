package fn

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

	fn.register("http.multipart", &FnEntry{
		Handler:     f.httpMultipart,
		Name:        "HTTP Multipart",
		Description: "Performs a multipart/form-data HTTP request with automatic file handling",
		Category:    f.category,
		Parameters: []FnParameter{
			{Name: "url", Description: "The URL to query", Mandatory: true},
			{Name: "method", Description: "HTTP method to use", Enum: []string{"POST", "PUT", "PATCH"}, Mandatory: true},
			{Name: "headers", Description: "Custom headers", Type: "object", Mandatory: false},
			{Name: "body", Description: "Fields and files for multipart", Type: "object", Mandatory: true},
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
			var j json.RawMessage
			if err := json.Unmarshal(response.Body(), &j); err != nil {
				return nil, fmt.Errorf("error unmarshaling JSON response: %w", err)
			}
			resp["body"] = j
		}

		if response.StatusCode() >= 400 {
			return nil, fmt.Errorf("HTTP request failed with status %d: %s", response.StatusCode(), string(response.Body()))
		}

		return utils.ReturnRaw(resp), nil
	})
}

type MultipartParams struct {
	Url     string            `json:"url" yaml:"url"`
	Method  string            `json:"method" yaml:"method"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    map[string]any    `json:"body" yaml:"body"` // key = field name, value = string, []byte, or base64 string
}

func (f *fnHttp) httpMultipart(j json.RawMessage) (json.RawMessage, error) {
	return utils.HandleJSON(j, func(params *MultipartParams) (json.RawMessage, error) {
		var bodyBuf bytes.Buffer
		writer := multipart.NewWriter(&bodyBuf)

		// Build multipart body
		for key, val := range params.Body {
			switch v := val.(type) {
			case string:
				if strings.HasPrefix(v, "data:") && strings.Contains(v, "base64,") {
					// Base64 string, decode and write as file
					parts := strings.SplitN(v, "base64,", 2)
					b, err := base64.StdEncoding.DecodeString(parts[1])
					if err != nil {
						return nil, fmt.Errorf("failed to decode base64 for field %s: %w", key, err)
					}
					part, err := writer.CreateFormFile(key, key)
					if err != nil {
						return nil, err
					}
					part.Write(b)
				} else {
					writer.WriteField(key, v)
				}
			case []byte:
				part, err := writer.CreateFormFile(key, key)
				if err != nil {
					return nil, err
				}
				part.Write(v)
			default:
				return nil, fmt.Errorf("unsupported value type for field %s", key)
			}
		}

		writer.Close() // finalize boundary

		// Prepare request
		req, err := http.NewRequest(params.Method, params.Url, &bodyBuf)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		for k, v := range params.Headers {
			req.Header.Set(k, v)
		}

		if _, ok := params.Headers["User-Agent"]; !ok {
			req.Header.Set("User-Agent", fmt.Sprintf("coda/%s", f.Fn.version))
		}

		// Execute
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Build response object
		result := map[string]any{
			"status":  resp.StatusCode,
			"headers": resp.Header,
			"body":    string(bodyBytes),
		}

		if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
			var jsonBody json.RawMessage
			if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
				result["body"] = jsonBody
			}
		}

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		return utils.ReturnRaw(result), nil
	})
}
