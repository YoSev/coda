package controller

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yosev/coda/internal/config"
	"github.com/yosev/coda/pkg/coda"
	"sigs.k8s.io/yaml"
)

func HandleYaml(c *gin.Context, payload []byte) {
	if c.Request.ContentLength > 0 {
		if c.Request.Header.Get("Content-Type") != "text/yaml" &&
			c.Request.Header.Get("Content-Type") != "application/x-yaml" &&
			c.Request.Header.Get("Content-Type") != "text/x-yaml" {
			c.JSON(400, gin.H{"error": "Content-Type must be text/yaml, application/x-yaml, or text/x-yaml"})
			return
		}

		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to read request body: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to read request body: %v\n", err)
			return
		}

		if payload != nil {
			payload, err = mergeYAML(payload, b)
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
	err := applyBlacklist(config.GetConfig().Blacklist, codaInstance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to apply blacklist: %v\n", err)
		os.Exit(1)
	}
	_, err = codaInstance.FromYaml(string(payload))
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to initiate coda from yaml: " + err.Error()})
		fmt.Fprintf(os.Stderr, "failed to initiate coda from yaml: %v\n", err)
		return
	}

	err = codaInstance.Run()
	addStatsToMetrics(codaInstance, err == nil)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to run coda: " + err.Error()})
		fmt.Fprintf(os.Stderr, "failed to run coda: %v\n", err)
		return
	}

	fmt.Printf("processed coda request with %d operations after %s\n", len(codaInstance.Operations), time.Since(start))
	codaInstance.Finish()
	y, err := yaml.Marshal(codaInstance)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to marshal coda to yaml: " + err.Error()})
		fmt.Fprintf(os.Stderr, "failed to marshal coda to yaml: %v\n", err)
		return
	}
	c.Header("Content-Type", "text/yaml")
	c.String(200, string(y))
}

func HandleYamlFile(c *gin.Context) {
	payload := downloadFile(c)
	HandleYaml(c, payload)
}

func mergeYAML(yamlA, yamlB []byte) ([]byte, error) {
	var mapA, mapB map[string]interface{}
	if err := yaml.Unmarshal(yamlA, &mapA); err != nil {
		return nil, fmt.Errorf("error unmarshalling yamlA: %w", err)
	}
	if err := yaml.Unmarshal(yamlB, &mapB); err != nil {
		return nil, fmt.Errorf("error unmarshalling yamlB: %w", err)
	}

	mergeMaps(mapA, mapB)

	mergedYAML, err := yaml.Marshal(mapA)
	if err != nil {
		return nil, fmt.Errorf("error marshalling merged map: %w", err)
	}

	return mergedYAML, nil
}
