package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	jsonSchema "github.com/yosev/coda/schema"
)

var schemaCmd = &cobra.Command{
	Use:                   "schema",
	DisableFlagsInUseLine: true,
	Short:                 "coda json schema",
	Run:                   schema,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}

func schema(cmd *cobra.Command, args []string) {
	fmt.Println(jsonSchema.GenerateSchema(version))
}
