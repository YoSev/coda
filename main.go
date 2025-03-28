package main

import (
	"os"

	_ "embed"

	"github.com/yosev/coda/cmd"
)

func main() {
	cmd.Execute(os.Args)
}
