package coda

import (
	"encoding/json"

	"github.com/yosev/coda/internal/coda/fn"
	"sigs.k8s.io/yaml"
)

// CodaSettings contains the settings for the coda engine
type CodaSettings struct {
	Debug bool `json:"debug" yaml:"debug"` // optional

	/**
	* The strict mode will:
	* - require store entries to be predefined
	 */
	Strict bool `json:"strict" yaml:"strict"` // optional
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
	Coda       CodaSettings               `json:"coda,omitempty" yaml:"coda,omitempty"`   // optional
	Store      map[string]json.RawMessage `json:"store,omitempty" yaml:"store,omitempty"` // optional
	Operations []Operation                `json:"operations" yaml:"operations"`           // mandatory

	fn     *fn.Fn
	source source
}

var defaultCodaSettings = CodaSettings{
	Debug: false,
}

func NewEmpty() *Coda {
	return new()
}

// NewFromJson creates a new Coda instance from a JSON string
//
//export NewFromJson
func NewFromJson(j string) (*Coda, error) {
	err := validateSchema(j)
	if err != nil {
		return nil, err
	}

	c := new()
	c.source = SOURCE_JSON // set source to JSON
	err = json.Unmarshal([]byte(j), c)
	if err != nil {
		return nil, err
	}

	c.debug("initialized new coda instance from json")
	return c, nil
}

// NewFromYaml creates a new Coda instance from a YAML string
func NewFromYaml(y string) (*Coda, error) {
	// TODO add schema validation for yaml
	c := new()
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
	return &Coda{fn: fn.New(), Store: make(map[string]json.RawMessage), Operations: []Operation{}, Coda: defaultCodaSettings}
}

// validateSchema validates the input against the JSON schema
func validateSchema(input string) error {
	// validate input against json-schema
	// schema := gojsonschema.NewStringLoader(schema.GetRaw())
	// jsonLoader := gojsonschema.NewStringLoader(input)
	// result, err := gojsonschema.Validate(schema, jsonLoader)

	// if err != nil {
	// 	return err
	// }
	// if !result.Valid() {
	// 	return fmt.Errorf("invalid JSON: %s", result.Errors())
	// }
	return nil
}
