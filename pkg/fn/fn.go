package fn

import (
	"encoding/json"
	"fmt"
)

type FnEntry struct {
	Handler     func(json.RawMessage) (json.RawMessage, error)
	Name        string
	Description string
	Category    FnCategory
	Parameters  []FnParameter
}

type FnParameter struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Mandatory   bool     `json:"mandatory" yaml:"mandatory"`
	Type        string   `json:"type" yaml:"type"`
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
}

type FnCategory string

const (
	FnCategoryFile      FnCategory = "File"
	FnCategoryString    FnCategory = "String"
	FnCategoryTime      FnCategory = "Time"
	FnCategoryIO        FnCategory = "I/O"
	FnCategoryMessaging FnCategory = "Messaging"
	FnCategoryOS        FnCategory = "OS"
	FnCategoryHTTP      FnCategory = "HTTP"
	FnCategoryHash      FnCategory = "Hash"
	FnCategoryMath      FnCategory = "Math"
	FnCategoryAI        FnCategory = "AI"
)

type FnHandler interface {
	Register()
}

type Fn struct {
	version string
	fns     map[string]*FnEntry
}

func (f *Fn) GetFns() map[string]*FnEntry {
	return f.fns
}

func New(version string) *Fn {
	f := &Fn{version: version, fns: make(map[string]*FnEntry)}

	// setup fn handlers
	(&FnHttp{Fn: f, Category: FnCategoryHTTP}).Register()
	(&FnAi{Fn: f, Category: FnCategoryAI}).Register()
	(&FnFile{Fn: f, Category: FnCategoryFile}).Register()
	(&FnS3{Fn: f, Category: FnCategoryFile}).Register()
	(&FnHash{Fn: f, Category: FnCategoryHash}).Register()
	(&FnIo{Fn: f, Category: FnCategoryIO}).Register()
	(&FnMath{Fn: f, Category: FnCategoryMath}).Register()
	(&FnMessage{Fn: f, Category: FnCategoryMessaging}).Register()
	(&FnOs{Fn: f, Category: FnCategoryOS}).Register()
	(&FnTime{Fn: f, Category: FnCategoryTime}).Register()
	(&FnString{Fn: f, Category: FnCategoryString}).Register()

	return f
}

func handleJSON[T any](j json.RawMessage, handler func(*T) (json.RawMessage, error)) (json.RawMessage, error) {
	var params T
	if err := json.Unmarshal(j, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters, invalid parameters: %v", err)
	}
	return handler(&params)
}

func returnRaw(obj any) json.RawMessage {
	b, _ := json.Marshal(obj)
	return json.RawMessage([]byte(b))
}

func addQuotes(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}
