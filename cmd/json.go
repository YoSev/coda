package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/yosev/coda/internal/coda"
)

var jsonCmd = &cobra.Command{
	Use:                   "j <coda script as json>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Short:                 "coda script as json",
	Run:                   jsonS,
}

var jsonFileCmd = &cobra.Command{
	Use:                   "jj <coda script as json file>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Example:               `coda jj script.coda.json`,
	Short:                 "coda script as json file",
	Run:                   jsonF,
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	rootCmd.AddCommand(jsonFileCmd)
}

func jsonS(cmd *cobra.Command, args []string) {
	input := args[0]
	if args[0] == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Printf("failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(b)
	}
	executeJ(input)
}

func jsonF(cmd *cobra.Command, args []string) {
	f, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("failed to read json file: %v\n", err)
		os.Exit(1)
	}
	script := string(f)
	executeJ(script)
}

func executeJ(j string) {
	c, err := coda.NewFromJson(j)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initiate coda from json file: %v\n", err)
		os.Exit(1)
	}

	err = c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
