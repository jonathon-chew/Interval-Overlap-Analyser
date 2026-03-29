# Interval-Overlap-Analyser

![Interval-Overlap-Analyser banner](./assets/logo.png)

## Summary

`Interval-Overlap-Analyser` is a range-analysis CLI for finding collisions, gaps, and merge opportunities across time windows such as meetings, maintenance periods, or booking slots. It makes interval algorithms concrete by giving them an obvious real-world use.

## Why This Project Exists

This project is meant to teach:

- interval sorting and sweep-line style reasoning,
- overlap detection and merge logic,
- gap discovery and free-slot calculation,
- how input shape affects algorithm choice and edge cases.

## Current Capabilities

- Read job lifecycle data from CSV.
- Group jobs whose `Start Time` falls within one minute of a window anchor.
- Build snapshots of jobs that have started but not yet dispatched.
- Exercise the parsing and grouping logic against large synthetic CSV datasets.

## Architecture Sketch

- A parser normalises incoming intervals into one internal representation.
- Sorting and sweep passes perform conflict and gap analysis.
- Output renderers present merged windows, conflicts, or free slots.
- Later versions can add calendar-style grouping by resource or person.

## Milestones

1. Parse intervals and detect pairwise overlaps.
2. Add merging and gap-finding over sorted ranges.
3. Add grouped analysis for multiple resources or calendars.
4. Add visual summaries, stricter validation, and performance tests on large inputs.

## Current Status

This project is an in-progress learning prototype rather than a finished CLI. The CSV parsing layer is implemented, there is active work on minute-based grouping and "started but not dispatched" windowing, and the repository includes synthetic datasets for stress-testing those ideas. The command-line interface is still minimal and the output is still debug-oriented.

## Development Notes

Useful commands during development:

- `go run ./cmd/Interval-Overlap-Analyser`
- `go build ./cmd/Interval-Overlap-Analyser`
- `go test ./...`

The current CLI entrypoint reads from `./testdata/fake_jobs.csv` by default.

## Project Structure

```text
cmd/Interval-Overlap-Analyser/  future interval-analysis entrypoint
internal/                       parsing, sorting, and overlap internals
pkg/                            optional reusable interval packages
doc/                            format notes and algorithm sketches
scripts/                        helper scripts
testdata/                       synthetic CSV datasets for local testing
```
