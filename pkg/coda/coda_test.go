package coda

import (
	"testing"
)

func TestNewFromJson_Valid(t *testing.T) {
	// minimal valid JSON with required operations field.
	input := `{
		"operations": []
	}`
	instance, err := New().FromJson(input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if instance == nil {
		t.Fatal("expected a valid instance, got nil")
	}
	// Check default settings
	if instance.Coda.Debug != false {
		t.Errorf("expected Coda.Debug to be false, got: %v", instance.Coda.Debug)
	}
	if instance.Store == nil {
		t.Error("expected Store to be initialized, but it is nil")
	}
	if instance.Operations == nil {
		t.Error("expected Operations to be initialized, but it is nil")
	}
}

func TestNewFromJson_Invalid(t *testing.T) {
	// Provide malformed JSON
	input := `{`
	_, err := New().FromJson(input)
	if err == nil {
		t.Fatal("expected an error for malformed JSON, but got nil")
	}
}

func TestNewFromYaml_Valid(t *testing.T) {
	// minimal valid YAML with settings and empty operations
	input := `
coda:
  debug: true
operations: []
`
	instance, err := New().FromYaml(input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if instance == nil {
		t.Fatal("expected a valid instance, got nil")
	}
	if instance.Coda.Debug != true {
		t.Errorf("expected Coda.Debug to be true, got: %v", instance.Coda.Debug)
	}
	if instance.Store == nil {
		t.Error("expected Store to be initialized, but it is nil")
	}
	if instance.Operations == nil {
		t.Error("expected Operations to be initialized, but it is nil")
	}
}

func TestNewFromYaml_Invalid(t *testing.T) {
	// Provide malformed YAML
	input := `:invalid_yaml`
	_, err := New().FromYaml(input)
	if err == nil {
		t.Fatal("expected an error for malformed YAML, but got nil")
	}
}
