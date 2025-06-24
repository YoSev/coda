package coda

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/yosev/coda/internal/fn"
	"sigs.k8s.io/yaml"

	_ "embed"
)

//go:embed .version
var VERSION string

// CodaSettings contains the settings for the coda engine
type CodaSettings struct {
	// Return logs of the coda run
	Logs bool `json:"logs" yaml:"logs"` // optional

	// Return runtime stats of the coda run
	Stats bool `json:"stats" yaml:"stats"` // optional

	// Return the coda settings and operations
	Extended bool `json:"extended" yaml:"extended"` // optional
}

// Operation is a single operation to be executed
type Operation struct {
	Entrypoint bool            `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"` // optional
	Action     string          `json:"action" yaml:"action"`                             // mandatory
	Params     json.RawMessage `json:"params,omitempty" yaml:"params,omitempty"`         // optional
	Store      string          `json:"store,omitempty" yaml:"store,omitempty"`           // optional
	OnSuccess  string          `json:"onSuccess,omitempty" yaml:"onSuccess,omitempty"`   // optional
	OnFail     string          `json:"onFail,omitempty" yaml:"onFail,omitempty"`         // optional
	Async      bool            `json:"async,omitempty" yaml:"async,omitempty"`           // optional
}

type source string

const (
	SOURCE_JSON source = "json" // source is JSON input
	SOURCE_YAML source = "yaml" // source is YAML input
)

// Coda is the main struct for the coda engine
type Coda struct {
	Coda       *CodaSettings              `json:"coda,omitempty" yaml:"coda,omitempty"`   // optional
	Logs       []string                   `json:"logs,omitempty" yaml:"logs,omitempty"`   // optional
	Stats      *CodaStats                 `json:"stats,omitempty" yaml:"stats,omitempty"` // optional
	Store      map[string]json.RawMessage `json:"store" yaml:"store"`
	Secrets    map[string]json.RawMessage `json:"secrets" yaml:"secrets"`
	Operations map[string]Operation       `json:"operations,omitempty" yaml:"operations,omitempty"` // mandatory

	fn        *fn.Fn
	source    source
	mutex     sync.RWMutex
	blacklist []OperationCategory `json:"-" yaml:"-"`
}

type codaDTO struct {
	Coda       *CodaSettings              `json:"coda,omitempty" yaml:"coda,omitempty"`
	Logs       []string                   `json:"logs,omitempty" yaml:"logs,omitempty"`
	Stats      *CodaStats                 `json:"stats,omitempty" yaml:"stats,omitempty"`
	Store      map[string]json.RawMessage `json:"store" yaml:"store"`
	Operations map[string]Operation       `json:"operations,omitempty" yaml:"operations,omitempty"`
}

// New creates a new Coda instance with default settings
func New() *Coda {
	return new()
}

// Run executes the coda operations and returns any error encountered
func (c *Coda) Run() error {
	return c.run()
}

func (c *Coda) ToDto() *codaDTO {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var out = &codaDTO{Store: c.Store}

	if c.Coda != nil {
		if c.Coda.Logs {
			out.Logs = c.Logs
		}
		if c.Coda.Stats {
			out.Stats = c.Stats
		}
		if c.Coda.Extended {
			out.Coda = c.Coda
			out.Operations = c.Operations
		}
	}

	return out
}

// Marshal the Coda instance to json or yaml based on the source
func (c *Coda) Marshal() ([]byte, error) {
	if c.source == SOURCE_YAML {
		// If the source is YAML, marshal to YAML
		return yaml.Marshal(c.ToDto())
	}
	return json.Marshal(c.ToDto())
}

// Blacklist categories of operations for this run
func (c *Coda) Blacklist(category OperationCategory) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.blacklist = append(c.blacklist, category)
}

// Create a new Coda instance from a JSON string
func (c *Coda) FromJson(j string) (*Coda, error) {
	err := c.validateSchema(j)
	if err != nil {
		return nil, err
	}

	c.source = SOURCE_JSON // set source to JSON
	err = json.Unmarshal([]byte(j), c)
	if err != nil {
		return nil, err
	}

	c.debug("initialized new coda instance from json")
	return c, nil
}

// Create a new Coda instance from a YAML string
func (c *Coda) FromYaml(y string) (*Coda, error) {
	if strings.HasPrefix(y, "{") || strings.HasPrefix(y, "[") {
		return nil, fmt.Errorf("input is not a valid YAML string")
	}

	// TODO add schema validation for yaml
	c.source = SOURCE_YAML // set source to JSON
	err := yaml.Unmarshal([]byte(y), c)
	if err != nil {
		return nil, err
	}

	c.debug("initialized new coda instance from yaml")
	return c, nil
}

// new creates a new Coda instance with default settings
func new() *Coda {
	c := &Coda{
		Coda: &CodaSettings{
			Logs:     false,
			Stats:    false,
			Extended: false,
		},
		Logs:       []string{},
		Store:      make(map[string]json.RawMessage),
		Operations: make(map[string]Operation),

		blacklist: []OperationCategory{},
		fn:        fn.New(VERSION),
		mutex:     sync.RWMutex{},
	}
	c.Stats = c.newStats()
	return c
}
