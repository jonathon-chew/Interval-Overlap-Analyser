package cli

import "github.com/jonathon-chew/Interval-Overlap-Analyser/internal/parse"

type Flags struct{}

func CLI(args []string) Flags {

	parse.StartToDispatch("./testdata/fake_jobs.csv")

	return Flags{}
}
