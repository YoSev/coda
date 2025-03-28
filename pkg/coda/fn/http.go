package fn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpReqParams struct {
	Url     string            `json:"url" yaml:"url"`
	Method  string            `json:"method" yaml:"method"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    json.RawMessage   `json:"body" yaml:"body"`
}

func (f *Fn) HttpReq(j json.RawMessage) (json.RawMessage, error) {

	return handleJSON(j, func(params *HttpReqParams) (json.RawMessage, error) {
		url, err := url.Parse(params.Url)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %s", err)
		}

		client := http.Client{Timeout: time.Duration(30) * time.Second}
		req := &http.Request{
			Method: params.Method,
			URL:    url,
			Header: http.Header{},
			Body:   io.NopCloser(strings.NewReader(string(params.Body))),
		}

		if string(params.Body) != "" {
			req.ContentLength = int64(len(params.Body))
		}

		for key, value := range params.Headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response body: %v", err)
		}

		responseData := map[string]interface{}{
			"statusCode": resp.StatusCode,
			"header":     resp.Header,
			"body":       string(bodyBytes),
		}

		if resp.Header.Get("Content-Type") == "application/json" {
			var b = make(map[string]interface{})
			err := json.Unmarshal(bodyBytes, &b)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal response body: %s", err)
			}
			responseData["body"] = b
		}

		if resp.StatusCode >= 400 {
			return returnRaw(responseData), fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
		}

		return returnRaw(responseData), nil
	})
}
