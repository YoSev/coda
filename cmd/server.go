package cmd

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/yosev/coda/internal/controller"
	"github.com/yosev/coda/internal/metrics"
	"github.com/yosev/coda/internal/middleware"
	"github.com/yosev/coda/internal/version"
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
		AllowCredentials: true,
		MaxAge:           86400 * time.Second,
	}))

	// setup version headers
	middleware.Version(router, version.VERSION)

	// setup basicAuth
	if basicAuth != nil && *basicAuth != "" {
		middleware.Auth(router, basicAuth)
	}

	// setup default route
	router.GET("/", func(c *gin.Context) {
		c.AbortWithStatus(200)
	})

	// setup json handler
	router.POST("/coda/j", func(c *gin.Context) {
		controller.HandleJson(c, blacklist, nil)
	})
	// setup json file handler
	router.POST("/coda/jj/*url", func(c *gin.Context) {
		controller.HandleJsonFile(c, blacklist)
	})
	// setup json file handler with custom payload
	router.POST("/coda/jjc/:key/*url", func(c *gin.Context) {
		controller.HandleJsonFile(c, blacklist)
	})

	// setup yaml handler
	router.POST("/coda/y", func(c *gin.Context) {
		controller.HandleYaml(c, blacklist, nil)
	})
	// setup yaml file handler
	router.POST("/coda/yy/*url", func(c *gin.Context) {
		controller.HandleYamlFile(c, blacklist)
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
