package main

import (
	"fmt"
	"os"

	purrpeek "github.com/nikhil25803/purrpeek/internal"
)

func main() {
	if err := purrpeek.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
