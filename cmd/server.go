package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/yosev/coda/internal/coda"
	"sigs.k8s.io/yaml"
)

var serverCmd = &cobra.Command{
	Use:                   "server",
	DisableFlagsInUseLine: true,
	Example:               `coda server`,
	Short:                 "coda server",
	Run:                   serverFn,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serverFn(cmd *cobra.Command, args []string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400 * time.Second,
	}))
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
		coda, err := coda.NewFromJson(string(b))
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to initiate coda from json: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to initiate coda from json: %v\n", err)
			return
		}
		err = coda.Run()
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to run coda: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to run coda: %v\n", err)
			return
		}

		fmt.Printf("processed coda request with %d operations after %s\n", len(coda.Operations), time.Since(start))
		c.JSON(200, coda)
	})
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
		coda, err := coda.NewFromYaml(string(b))
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to initiate coda from yaml: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to initiate coda from yaml: %v\n", err)
			return
		}
		err = coda.Run()
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to run coda: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to run coda: %v\n", err)
			return
		}

		fmt.Printf("processed coda request with %d operations after %s\n", len(coda.Operations), time.Since(start))
		y, err := yaml.Marshal(coda)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to marshal coda to yaml: " + err.Error()})
			fmt.Fprintf(os.Stderr, "failed to marshal coda to yaml: %v\n", err)
			return
		}
		c.Header("Content-Type", "text/yaml")
		c.String(200, string(y))
	})

	fmt.Printf("coda server started on port %d\n", *port)
	router.Run(fmt.Sprintf(":%d", *port))
}
