package coda

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"github.com/yosev/coda/pkg/coda/fn"
	"sigs.k8s.io/yaml"
)

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
	Action string          `json:"action" yaml:"action"`                     // mandatory
	Params json.RawMessage `json:"params,omitempty" yaml:"params,omitempty"` // optional
	Store  string          `json:"store,omitempty" yaml:"store,omitempty"`   // optional
	OnFail []Operation     `json:"onFail,omitempty" yaml:"onFail,omitempty"` // optional
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
	Operations []Operation                `json:"operations,omitempty" yaml:"operations,omitempty"` // mandatory

	fn        *fn.Fn
	source    source
	blacklist []OperationCategory `json:"-" yaml:"-"`
}

func New() *Coda {
	return new()
}

func (c *Coda) Run() error {
	return c.run()
}

// Finish the Coda instance by applying the coda.Coda settings
func (c *Coda) Finish() {
	if !c.Coda.Stats {
		c.Stats = nil
	}
	if !c.Coda.Logs {
		c.Logs = nil
	}
	if !c.Coda.Extended {
		c.Coda = nil
		c.Operations = nil
	}
}

// Blacklist categories of operations for this run
func (c *Coda) Blacklist(category OperationCategory) {
	c.blacklist = append(c.blacklist, category)
}

// NewFromJson creates a new Coda instance from a JSON string
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

// NewFromYaml creates a new Coda instance from a YAML string
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
		Operations: []Operation{},

		blacklist: []OperationCategory{},
		fn:        fn.New(),
	}
	c.Stats = c.newStats()
	return c
}

// validateSchema validates the input against the JSON schema
func (c *Coda) validateSchema(input string) error {
	// TODO allow string wildcards for type integer (type: ["integer","string"])
	schema := gojsonschema.NewStringLoader(c.Schema())
	jsonLoader := gojsonschema.NewStringLoader(input)
	result, err := gojsonschema.Validate(schema, jsonLoader)

	if err != nil {
		return err
	}
	if !result.Valid() {
		return fmt.Errorf("invalid JSON: %s", result.Errors())
	}
	return nil
}
