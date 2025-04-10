package schema

import (
	_ "embed"
	"encoding/json"

	"github.com/yosev/coda/pkg/coda"
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
	Type  string   `json:"type,omitempty"`
	Enum  []string `json:"enum,omitempty"`
	Const string   `json:"const,omitempty"`
}

// GenerateSchema loads and populates the JSON schema
func GenerateSchema(version string) string {
	err := json.Unmarshal([]byte(jsonSchemaRaw), jsonSchema)
	if err != nil {
		panic("Failed to parse JSON schema: " + err.Error())
	}
	jsonSchema.populateSchema(version)

	b, _ := json.Marshal(jsonSchema)
	return string(b)
}

func GetRawSchema() string {
	return string(jsonSchemaRaw)
}

func (s *Schema) populateSchema(version string) {
	s.Version = version
	s.Defs = map[string]map[string]interface{}{}
	s.Properties.Operations = &SchemaOperationsProperty{
		Type: "array",
	}

	// Collect all specific operation definitions
	unifiedAnyOf := []map[string]interface{}{}

	for _, operation := range coda.NewEmpty().GetOperations() {
		paramDefinitions := map[string]SchemaOperationParams{}
		requiredParamNames := []string{}

		for _, parameter := range operation.Parameters {
			param := SchemaOperationParams{Type: "string"}
			if parameter.Type != "" && parameter.Type != "any" {
				param.Type = parameter.Type
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
			},
			"required":             requiredFields,
			"additionalProperties": false,
		}

		if len(paramDefinitions) > 0 {
			paramProps := map[string]interface{}{}
			for k, v := range paramDefinitions {
				prop := map[string]interface{}{}
				if v.Type != "" {
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
