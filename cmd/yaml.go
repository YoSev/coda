package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yosev/coda/pkg/coda"
	"gopkg.in/yaml.v3"
)

var yamlCmd = &cobra.Command{
	Use:                   "y <coda script as yaml>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Short:                 "coda script as yaml",
	Run:                   yamlS,
}

var yamlFileCmd = &cobra.Command{
	Use:                   "yy <coda script as yaml file>",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Example:               `coda y script.coda.yaml`,
	Short:                 "coda script as yaml file",
	Run:                   yamlF,
}

func init() {
	rootCmd.AddCommand(yamlCmd)
	rootCmd.AddCommand(yamlFileCmd)
}

func yamlS(cmd *cobra.Command, args []string) {
	input := args[0]
	if args[0] == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
		input = string(b)
	}
	executeY(input)
}

func yamlF(cmd *cobra.Command, args []string) {
	f, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read yaml file: %v\n", err)
		os.Exit(1)
	}
	script := string(f)

	// remove shebang if present
	if strings.HasPrefix(script, "#!") {
		script = strings.SplitN(script, "\n", 2)[1]
	}

	executeY(script)
}

func executeY(j string) {
	c, err := coda.New().FromYaml(j)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initiate coda from yaml file: %v\n", err)
		os.Exit(1)
	}

	err = c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	} else {
		c.CleanUp()

		// Convert to map first to handle json.RawMessage
		var tmpMap map[string]interface{}
		jsonData, err := json.Marshal(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal coda to json: %v\n", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(jsonData, &tmpMap); err != nil {
			fmt.Fprintf(os.Stderr, "failed to unmarshal json to map: %v\n", err)
			os.Exit(1)
		}

		b, err := yaml.Marshal(tmpMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal coda response: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s\n", b)
	}
}
