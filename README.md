# Interval-Overlap-Analyser

## Summary

`Interval-Overlap-Analyser` is a range-analysis CLI for finding collisions, gaps, and merge opportunities across time windows such as meetings, maintenance periods, or booking slots. It makes interval algorithms concrete by giving them an obvious real-world use.

## Why This Project Exists

This project is meant to teach:

- interval sorting and sweep-line style reasoning,
- overlap detection and merge logic,
- gap discovery and free-slot calculation,
- how input shape affects algorithm choice and edge cases.

## Planned Capabilities

- Read intervals from CSV, JSON, or simple text input.
- Detect overlaps and report conflicting ranges clearly.
- Merge intersecting ranges into compact summaries.
- Find gaps or available windows between busy periods.

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

This project is currently scaffolded but not implemented. The folder layout is ready, but the parsing, validation, and interval-analysis logic still need to be built.

## Development Notes

Planned commands once implementation begins:

- `go run ./cmd/Interval-Overlap-Analyser`
- `go build ./cmd/Interval-Overlap-Analyser`
- `go test ./...`

## Project Structure

```text
cmd/Interval-Overlap-Analyser/  future interval-analysis entrypoint
internal/                       parsing, sorting, and overlap internals
pkg/                            optional reusable interval packages
doc/                            format notes and algorithm sketches
scripts/                        helper scripts
```
