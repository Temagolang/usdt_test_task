package main

import (
	"os"

	"github.com/example/grinex-rates-service/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
