package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var port *int

var version = ""

var rootCmd = &cobra.Command{
	Use:   "coda",
	Args:  cobra.ExactArgs(1),
	Short: "coda is an interpreter for coda scripts",
	Long:  "coda is an interpreter for coda scripts which come in form of coda, json or yaml",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute(args []string, v string) {
	version = v
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.ParseFlags(args)

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		port = serverCmd.PersistentFlags().IntP("port", "p", 3000, "port to run the server on")
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
