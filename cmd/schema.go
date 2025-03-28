package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yosev/coda/pkg/coda"
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
	fmt.Println(coda.New().Schema())
}
