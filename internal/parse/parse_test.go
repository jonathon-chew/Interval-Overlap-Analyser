package parse

import "testing"

func TestBuildStartToDispatchWindowsSnapshotsBeforeRemoval(t *testing.T) {
	entries := []Entry{
		{
			Job_Number:    "1001",
			Start_Time:    parseTime("2026-03-24 09:00:00"),
			Dispatch_Time: parseTime("2026-03-24 09:10:00"),
		},
		{
			Job_Number:    "1002",
			Start_Time:    parseTime("2026-03-24 09:05:00"),
			Dispatch_Time: parseTime("2026-03-24 09:20:00"),
		},
		{
			Job_Number:    "1003",
			Start_Time:    parseTime("2026-03-24 09:11:00"),
			Dispatch_Time: parseTime("2026-03-24 09:30:00"),
		},
	}

	windows := buildStartToDispatchWindows(entries)
	if len(windows.Windows) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(windows.Windows))
	}

	first := windows.Windows[0]
	if len(first.Entries) != 2 {
		t.Fatalf("expected first window to have 2 entries, got %d", len(first.Entries))
	}
	if first.Entries[0].Job_Number != "1001" || first.Entries[1].Job_Number != "1002" {
		t.Fatalf("unexpected first window job numbers: %+v", first.Entries)
	}
	if got := first.EndTime.Format("2006-01-02 15:04:05"); got != "2026-03-24 09:11:00" {
		t.Fatalf("expected first window end time to be 2026-03-24 09:11:00, got %s", got)
	}
}

func TestBuildStartToDispatchWindowsClonesSnapshots(t *testing.T) {
	entries := []Entry{
		{
			Job_Number:    "2001",
			Start_Time:    parseTime("2026-03-24 10:00:00"),
			Dispatch_Time: parseTime("2026-03-24 10:05:00"),
		},
		{
			Job_Number:    "2002",
			Start_Time:    parseTime("2026-03-24 10:01:00"),
			Dispatch_Time: parseTime("2026-03-24 10:08:00"),
		},
		{
			Job_Number:    "2003",
			Start_Time:    parseTime("2026-03-24 10:06:00"),
			Dispatch_Time: parseTime("2026-03-24 10:10:00"),
		},
	}

	windows := buildStartToDispatchWindows(entries)
	if len(windows.Windows) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(windows.Windows))
	}

	first := windows.Windows[0]
	second := windows.Windows[1]

	if len(first.Entries) != 2 || first.Entries[0].Job_Number != "2001" || first.Entries[1].Job_Number != "2002" {
		t.Fatalf("first snapshot was mutated: %+v", first.Entries)
	}
	if len(second.Entries) != 2 || second.Entries[0].Job_Number != "2002" || second.Entries[1].Job_Number != "2003" {
		t.Fatalf("unexpected second snapshot: %+v", second.Entries)
	}
}
