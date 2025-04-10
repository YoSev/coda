package coda

import (
	"fmt"
	"os"
)

func (c *Coda) debug(msg string) {
	if c.Coda.Debug {
		fmt.Fprint(os.Stderr, "debug: "+msg+"\n")
	}
}
