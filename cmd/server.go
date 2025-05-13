package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/yosev/coda/pkg/coda"
	"github.com/yosev/coda/pkg/metrics"
	"sigs.k8s.io/yaml"
)

var serverCmd = &cobra.Command{
	Use:                   "server",
	DisableFlagsInUseLine: true,
	Example:               `coda server`,
	Short:                 "coda server",
	Run:                   serverFn,
}

var port *int
var blacklist *[]string
var basicAuth *string

func init() {
	rootCmd.AddCommand(serverCmd)
	port = serverCmd.PersistentFlags().IntP("port", "p", 3000, "port to run the server on")
	blacklist = serverCmd.PersistentFlags().StringSliceP("blacklist", "b", []string{}, "comma separated list of blacklisted operation categories")
	basicAuth = serverCmd.PersistentFlags().StringP("auth", "a", "", "base64 encoded username:password for basic auth")
}

func serverFn(cmd *cobra.Command, args []string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// setup cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400 * time.Second,
	}))

	// setup basicAuth
	if basicAuth != nil && *basicAuth != "" {
		router.Use(func(c *gin.Context) {
			// do not apply authorization for GET requests (eg. for health check)
			if c.Request.Method == "GET" && c.Request.URL.Path == "/" {
				c.Next()
				return
			}

			auth := c.Request.Header.Get("Authorization")
			if auth == "" {
				c.Header("WWW-Authenticate", "Basic realm=\"coda\"")
				c.AbortWithStatus(401)
				return
			}
			if auth != "Basic "+*basicAuth {
				c.AbortWithStatus(401)
				return
			}
		})
	}

	// setup default route
	router.GET("/", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})

	// setup json handler
	router.POST("/coda/j", func(c *gin.Context) {
		start := time.Now()
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

		codaInstance := coda.New()
		err = applyBlacklist(codaInstance)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to apply blacklist: %v\n", err)
			os.Exit(1)
		}
		_, err = codaInstance.FromJson(string(b))
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
	})

	// setup yaml handler
	router.POST("/coda/y", func(c *gin.Context) {
		start := time.Now()
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

		codaInstance := coda.New()
		err = applyBlacklist(codaInstance)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to apply blacklist: %v\n", err)
			os.Exit(1)
		}
		_, err = codaInstance.FromYaml(string(b))
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to initiate coda from yaml: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to initiate coda from yaml: %v\n", err)
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
		y, err := yaml.Marshal(codaInstance)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to marshal coda to yaml: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to marshal coda to yaml: %v\n", err)
			return
		}
		c.Header("Content-Type", "text/yaml")
		c.String(200, string(y))
	})

	// setup metrics
	m := metrics.Registry()
	h := promhttp.HandlerFor(m, promhttp.HandlerOpts{EnableOpenMetrics: true})
	router.GET("/metrics", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	fmt.Printf("coda server started on port %d\n", *port)
	router.Run(fmt.Sprintf(":%d", *port))
}

func applyBlacklist(codaInstance *coda.Coda) error {
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
