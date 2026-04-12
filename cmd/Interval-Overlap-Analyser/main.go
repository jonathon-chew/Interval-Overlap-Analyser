package main

import (
	"os"

	"github.com/jonathon-chew/Interval-Overlap-Analyser/internal/cli"
)

func main() {
	_ = cli.CLI(os.Args[1:])
}
