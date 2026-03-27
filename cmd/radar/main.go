package main

import (
	"os"

	"github.com/infrashift/technology-radar/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
