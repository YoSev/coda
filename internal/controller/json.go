package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yosev/coda/pkg/coda"
)

func HandleJson(c *gin.Context, blacklist *[]string, payload []byte) {
	if c.Request.ContentLength > 0 {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(400, gin.H{"error": "Content-Type must be application/json"})
			return
		}

		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to read request body: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to read request body: %v\n", err)
			return
		}

		if c.Param("key") != "" {
			key := c.Param("key")
			payload := &coda.Coda{
				Store: map[string]json.RawMessage{key: b},
			}
			b, err = json.Marshal(payload)
			if err != nil {
				c.JSON(400, gin.H{"error": "failed to inject request body into store: " + err.Error()})
				fmt.Fprintf(os.Stderr, "failed to inject request body into store: %v\n", err)
				return
			}
		}

		if payload != nil {
			payload, err = mergeJson(payload, b)
			if err != nil {
				c.JSON(400, gin.H{"error": "failed to merge payloads"})
				return
			}
		} else {
			payload = b
		}
	}

	start := time.Now()
	codaInstance := coda.New()
	err := applyBlacklist(blacklist, codaInstance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to apply blacklist: %v\n", err)
		os.Exit(1)
	}
	_, err = codaInstance.FromJson(string(payload))
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to initiate coda from json: " + err.Error()})
		fmt.Fprintf(os.Stderr, "failed to initiate coda from json: %v\n", err)
		return
	}

	err = codaInstance.Run()
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to run coda: " + err.Error()})
		fmt.Fprintf(os.Stderr, "failed to run coda: %v\n", err)
		return
	}

	fmt.Printf("processed coda request with %d operations after %s\n", len(codaInstance.Operations), time.Since(start))
	// remove coda and operations from the response as they are not needed
	codaInstance.Coda = nil
	codaInstance.Operations = nil
	c.JSON(200, codaInstance)
}

func HandleJsonFile(c *gin.Context, blacklist *[]string) {
	payload := downloadFile(c)
	HandleJson(c, blacklist, payload)
}

func mergeJson(jsonA, jsonB []byte) ([]byte, error) {
	var mapA, mapB map[string]interface{}
	if err := json.Unmarshal(jsonA, &mapA); err != nil {
		return nil, fmt.Errorf("error unmarshalling jsonA: %w", err)
	}
	if err := json.Unmarshal(jsonB, &mapB); err != nil {
		return nil, fmt.Errorf("error unmarshalling jsonB: %w", err)
	}

	mergeMaps(mapA, mapB)

	mergedJSON, err := json.Marshal(mapA)
	if err != nil {
		return nil, fmt.Errorf("error marshalling merged map: %w", err)
	}

	return mergedJSON, nil
}
