package coda

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"github.com/yosev/coda/pkg/coda/fn"
	"sigs.k8s.io/yaml"
)

// CodaSettings contains the settings for the coda engine
type CodaSettings struct {
	Debug bool `json:"debug" yaml:"debug"` // optional
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
	Coda       *CodaSettings              `json:"coda,omitempty" yaml:"coda,omitempty"` // optional
	Logs       []string                   `json:"logs,omitempty" yaml:"logs,omitempty"` // optional
	Store      map[string]json.RawMessage `json:"store" yaml:"store"`
	Operations []Operation                `json:"operations,omitempty" yaml:"operations,omitempty"` // mandatory

	fn        *fn.Fn
	source    source
	Blacklist []OperationCategory `json:"-" yaml:"-"`
}

func New() *Coda {
	return new()
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
	return &Coda{
		Coda: &CodaSettings{
			Debug: false,
		},
		Logs:       []string{},
		Store:      make(map[string]json.RawMessage),
		Operations: []Operation{},

		Blacklist: []OperationCategory{},
		fn:        fn.New(),
	}
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
