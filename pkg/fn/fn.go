package fn

import (
	"encoding/json"
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

type fnHandler interface {
	init(fn *Fn)
}

type Fn struct {
	version string
	fns     map[string]*FnEntry
}

func (f *Fn) GetFns() map[string]*FnEntry {
	return f.fns
}

func (f *Fn) register(name string, entry *FnEntry) {
	if _, exists := f.fns[name]; exists {
		panic("function already registered: " + name)
	}
	f.fns[name] = entry
}

func New(version string) *Fn {
	f := &Fn{version: version, fns: make(map[string]*FnEntry)}

	// setup fn handlers
	var h = []fnHandler{
		&fnHttp{category: FnCategoryAI},
		&fnAi{category: FnCategoryAI},
		&fnFile{category: FnCategoryFile},
		&fnS3{category: FnCategoryFile},
		&fnHash{category: FnCategoryHash},
		&fnIo{category: FnCategoryIO},
		&fnMath{category: FnCategoryMath},
		&fnMessage{category: FnCategoryMessaging},
		&fnOs{category: FnCategoryOS},
		&fnTime{category: FnCategoryTime},
		&fnString{category: FnCategoryString},
	}
	for _, handler := range h {
		handler.init(f)
	}

	return f
}
