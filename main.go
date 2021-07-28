package main

import (
	"fmt"
	"github.com/truewhitespace/key-rotation/cmd"
	"os"
)

func main() {
	root := cmd.NewRoot()
	err := root.Execute()
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error()); err != nil {
			panic(err)
		}
		os.Exit(-1)
	}
}
