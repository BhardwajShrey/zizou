package main

import (
	"os"

	"github.com/shreybhardwaj/zizou/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
