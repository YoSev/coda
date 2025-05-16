package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yosev/coda/internal/config"
	"github.com/yosev/coda/internal/metrics"
	"github.com/yosev/coda/pkg/coda"
)

func addStatsToMetrics(c *coda.Coda, success bool) {
	metrics.Inc("coda_total")

	if success {
		metrics.Inc("coda_successful_total")
	} else {
		metrics.Inc("coda_failed_total")
	}

	metrics.IncValue("coda_runtime_total", c.Stats.CodaRuntimeTotalMs)
	metrics.IncValue("operations_runtime_total", c.Stats.OperationsRuntimeTotalMs)
	metrics.IncValue("operations_total", c.Stats.OperationsTotal)
	metrics.IncValue("operations_successful_total", c.Stats.OperationsSuccessfulTotal)
	metrics.IncValue("operations_failed_total", c.Stats.OperationsFailedTotal)
	metrics.IncValue("operations_blacklisted_total", c.Stats.OperationsBlacklistedTotal)
	metrics.IncValue("variables_total", c.Stats.VariablesTotal)
	metrics.IncValue("variables_failed_total", c.Stats.VariablesFailedTotal)
	metrics.IncValue("variables_successful_total", c.Stats.VariablesSuccessfulTotal)

	if config.GetConfig().InfluxDB != nil && *config.GetConfig().InfluxDB != "" && c.Stats != nil {
		err := metrics.SendStatsToInfluxDB(*config.GetConfig().InfluxDB, c.Stats)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to send stats to InfluxDB: %v", err)
		}
	}
}

func downloadFile(c *gin.Context) []byte {
	url, err := url.Parse(strings.TrimLeft(c.Param("url"), "/"))

	if err != nil {
		c.JSON(400, gin.H{"error": "the given url is not valid"})
		return nil
	}

	const maxFileSize = 1 * 1024 * 1024 // 1 MB

	var payload []byte

	switch url.Scheme {
	case "http", "https":
		resp, err := http.Get(url.String())
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to download given file"})
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode > 400 {
			c.JSON(400, gin.H{"error": "failed to download given file"})
			return nil
		}

		if resp.ContentLength > maxFileSize {
			c.JSON(400, gin.H{"error": "file size exceeds the maximum allowed limit"})
			return nil
		}

		payload, err = io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to read response body"})
			return nil
		}
	default:
		c.JSON(400, gin.H{"error": "the given url is not valid"})
		return nil
	}

	return payload
}

func applyBlacklist(blacklist *[]string, c *coda.Coda) error {
	if *blacklist == nil || len(*blacklist) == 0 {
		return nil
	}

	for _, category := range *blacklist {
		switch category {
		case "file":
			c.Blacklist(coda.OperationCategoryFile)
		case "string":
			c.Blacklist(coda.OperationCategoryString)
		case "time":
			c.Blacklist(coda.OperationCategoryTime)
		case "io":
			c.Blacklist(coda.OperationCategoryIO)
		case "os":
			c.Blacklist(coda.OperationCategoryOS)
		case "http":
			c.Blacklist(coda.OperationCategoryHTTP)
		case "hash":
			c.Blacklist(coda.OperationCategoryHash)
		case "math":
			c.Blacklist(coda.OperationCategoryMath)
		default:
			return fmt.Errorf("unknown blacklist category: %s", category)
		}
	}

	return nil
}

func mergeMaps(dst, src map[string]interface{}) {
	for key, srcValue := range src {
		if dstValue, exists := dst[key]; exists {
			if dstMap, ok1 := dstValue.(map[string]interface{}); ok1 {
				if srcMap, ok2 := srcValue.(map[string]interface{}); ok2 {
					mergeMaps(dstMap, srcMap)
					continue
				}
			}
		}
		dst[key] = srcValue
	}
}
