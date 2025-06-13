package coda

import (
	"fmt"
	"testing"

	_ "embed"
)

//go:embed _test/assets/coda.test.json
var testCodaJSON string

func TestNew(t *testing.T) {
	var c *Coda
	var err error
	var result []byte

	c = New()

	_, err = c.FromJson(testCodaJSON)
	if err != nil {
		t.Fatalf("failed to load coda from JSON: %v", err)
	}

	err = c.Run()
	if err != nil {
		t.Fatalf("failed to run coda: %v", err)
	}

	result, err = c.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal coda: %v", err)
	}

	fmt.Println(string(result))
}
