package main

import (
	"os"

	"Interval-Overlap-Analyser/internal/cli"
)

func main() {
	_ = cli.CLI(os.Args[1:])
}
