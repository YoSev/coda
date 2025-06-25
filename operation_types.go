package coda

import "encoding/json"

type OperationParameter struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Mandatory   bool     `json:"mandatory" yaml:"mandatory"`
	Type        string   `json:"type" yaml:"type"`
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
}

type OperationHandler struct {
	Fn          func(*Coda, json.RawMessage) (json.RawMessage, error)
	Name        string
	Description string
	Category    OperationCategory
	Parameters  []OperationParameter
}

type OperationCategory string

const (
	OperationCategoryFile      OperationCategory = "File"
	OperationCategoryString    OperationCategory = "String"
	OperationCategoryTime      OperationCategory = "Time"
	OperationCategoryIO        OperationCategory = "I/O"
	OperationCategoryMessaging OperationCategory = "Messaging"
	OperationCategoryOS        OperationCategory = "OS"
	OperationCategoryHTTP      OperationCategory = "HTTP"
	OperationCategoryHash      OperationCategory = "Hash"
	OperationCategoryMath      OperationCategory = "Math"
	OperationCategoryAI        OperationCategory = "AI"
)
