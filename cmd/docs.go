package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yosev/coda/pkg/coda"
)

var docsCmd = &cobra.Command{
	Use:                   "docs",
	DisableFlagsInUseLine: true,
	Short:                 "coda operations documentation",
	Run:                   docs,
}

func init() {
	rootCmd.AddCommand(docsCmd)
}

func docs(cmd *cobra.Command, args []string) {
	c := coda.NewEmpty()
	for _, operation := range c.GetOperations() {
		fmt.Printf("Action: %s\nDescription: %s\nCategory: %s\n", operation.Name, operation.Description, operation.Category)
		if len(operation.Parameters) > 0 {
			fmt.Printf("Parameters:\n")
		}
		for _, v := range operation.Parameters {
			name := v.Name
			if v.Mandatory {
				name = "(*)" + name
			}
			fmt.Printf("  - %s: %s\n", name, v.Description)
		}
		fmt.Printf("\n")
	}
}
