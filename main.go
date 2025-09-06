package main

import (
	"fmt"
	"os"

	"github.com/luizvilasboas/commit-hooks/internal/config"
	"github.com/luizvilasboas/commit-hooks/internal/tui"
)

func main() {
	var initialMsg string
	if len(os.Args) > 1 {
		initialMsg = os.Args[1]
	}

	cfg := config.Load()

	output, err := tui.Run(initialMsg, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}

	if output != "" {
		fmt.Print(output)
	}
}
