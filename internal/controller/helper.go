package controller

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yosev/coda/pkg/coda"
)

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

func applyBlacklist(blacklist *[]string, codaInstance *coda.Coda) error {
	if len(*blacklist) == 0 {
		return nil
	}

	for _, category := range *blacklist {
		switch category {
		case "file":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryFile)
		case "string":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryString)
		case "time":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryTime)
		case "io":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryIO)
		case "os":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryOS)
		case "http":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryHTTP)
		case "hash":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryHash)
		case "math":
			codaInstance.Blacklist = append(codaInstance.Blacklist, coda.OperationCategoryMath)
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
