package coda

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed coda.schema.json
var jsonSchemaRaw string
var jsonSchema = &Schema{}

type Schema struct {
	Schema               string                            `json:"$schema"`
	Title                string                            `json:"title,omitempty"`
	Version              string                            `json:"version,omitempty"`
	Type                 string                            `json:"type"`
	Properties           SchemaProperties                  `json:"properties"`
	Required             []string                          `json:"required,omitempty"`
	AdditionalProperties bool                              `json:"additionalProperties"`
	Defs                 map[string]map[string]interface{} `json:"$defs,omitempty"`
}

type SchemaProperties struct {
	Coda       *SchemaCodaProperty       `json:"coda,omitempty"`
	Store      *SchemaStoreProperty      `json:"store,omitempty"`
	Operations *SchemaOperationsProperty `json:"operations,omitempty"`
}

type SchemaCodaProperty struct {
	Type                 string                        `json:"type"`
	Properties           map[string]SchemaCodaProperty `json:"properties,omitempty"`
	Required             []string                      `json:"required,omitempty"`
	AdditionalProperties bool                          `json:"additionalProperties"`
}

type SchemaStoreProperty struct {
	Type                 string   `json:"type"`
	Required             []string `json:"required,omitempty"`
	AdditionalProperties bool     `json:"additionalProperties"`
}

type SchemaOperationsProperty struct {
	Type  string            `json:"type"`
	Items *SchemaItemsField `json:"items"`
}

type SchemaItemsField struct {
	Ref string `json:"$ref"`
}

type SchemaOperationAction struct {
	Type  string   `json:"type"`
	Enum  []string `json:"enum,omitempty"`
	Const string   `json:"const,omitempty"`
}

type SchemaOperationStore struct {
	Type string `json:"type"`
}

type SchemaOperationParamsWrapper struct {
	Properties           map[string]SchemaOperationParams `json:"properties"`
	Required             []string                         `json:"required,omitempty"`
	AdditionalProperties bool                             `json:"additionalProperties"`
}

type SchemaOperationParams struct {
	Type  []string `json:"type,omitempty"`
	Enum  []string `json:"enum,omitempty"`
	Const string   `json:"const,omitempty"`
}

func init() {
	new().Schema()
}

var schema = ""

// Get the JSON schema for the Coda engine.
func (c *Coda) Schema() string {
	if schema == "" {
		err := json.Unmarshal([]byte(jsonSchemaRaw), jsonSchema)
		if err != nil {
			panic("Failed to parse JSON schema: " + err.Error())
		}
		jsonSchema.populateSchema(VERSION)

		b, _ := json.Marshal(jsonSchema)
		schema = string(b)
	}
	return schema
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

// populateSchema fills the Schema struct with operation definitions and properties.
func (s *Schema) populateSchema(version string) {
	s.Version = version
	s.Defs = map[string]map[string]interface{}{}
	s.Properties.Operations = &SchemaOperationsProperty{
		Type: "array",
	}

	// Collect all specific operation definitions
	unifiedAnyOf := []map[string]interface{}{}

	for _, operation := range New().GetOperations() {
		paramDefinitions := map[string]SchemaOperationParams{}
		requiredParamNames := []string{}

		for _, parameter := range operation.Parameters {
			param := SchemaOperationParams{Type: []string{"string"}}
			if parameter.Type != "" {
				if parameter.Type == "any" {
					param.Type = []string{}
				} else if len(strings.Split(parameter.Type, ",")) > 1 {
					param.Type = strings.Split(parameter.Type, ",")
				} else {
					param.Type = []string{parameter.Type}
				}
			}
			if len(param.Type) > 0 && !slices.Contains(param.Type, "string") {
				param.Type = append(param.Type, "string") // Ensure string is always included to support $wildcards

			}
			if parameter.Enum != nil {
				param.Enum = parameter.Enum
			}
			paramDefinitions[parameter.Name] = param
			if parameter.Mandatory {
				requiredParamNames = append(requiredParamNames, parameter.Name)
			}
		}

		requiredFields := []string{"action"}
		if len(requiredParamNames) > 0 {
			requiredFields = append(requiredFields, "params")
		}

		defName := "Operation_" + operation.Name

		opSchema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":  "string",
					"const": operation.Name,
				},
				"store": map[string]interface{}{
					"type": "string",
				},
				"onFail": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/$defs/Operation",
					},
				},
				"async": map[string]interface{}{
					"type": "boolean",
				},
			},
			"required":             requiredFields,
			"additionalProperties": false,
		}

		if len(paramDefinitions) > 0 {
			paramProps := map[string]interface{}{}
			for k, v := range paramDefinitions {
				prop := map[string]interface{}{}
				if len(v.Type) != 0 {
					prop["type"] = v.Type
				}
				if len(v.Enum) > 0 {
					prop["enum"] = v.Enum
				}
				if v.Const != "" {
					prop["const"] = v.Const
				}
				paramProps[k] = prop
			}

			opSchema["properties"].(map[string]interface{})["params"] = map[string]interface{}{
				"type":                 "object",
				"properties":           paramProps,
				"required":             requiredParamNames,
				"additionalProperties": false,
			}
		}

		// Add operation to $defs
		s.Defs[defName] = opSchema

		// Add ref to unified Operation.anyOf
		unifiedAnyOf = append(unifiedAnyOf, map[string]interface{}{
			"$ref": "#/$defs/" + defName,
		})
	}

	// Final unified Operation definition
	s.Defs["Operation"] = map[string]interface{}{
		"anyOf": unifiedAnyOf,
	}

	// Assign operations.items to $ref Operation
	s.Properties.Operations.Items = &SchemaItemsField{
		Ref: "#/$defs/Operation",
	}
}
