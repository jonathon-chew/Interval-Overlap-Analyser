package parse

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type Entry struct {
	Job_Number       string
	Job_Type         string
	Completed_By     string
	Start_Time       time.Time
	Dispatch_Time    time.Time
	In_Progress_Time time.Time
	Complete_Time    time.Time
}

type Windows struct {
	Windows []Window
}

type Window struct {
	Entries     []Entry
	Job_Numbers []string
	StartTime   time.Time
	EndTime     time.Time
}

func parseTime(entry string) time.Time {
	value, err := time.Parse("2006-01-02 15:04:05", string(entry))
	if err != nil {
		return time.Time{}
	}

	return value
}

func readEntries(fileName string) ([]Entry, error) {
	filePointer, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	defer filePointer.Close()

	scanner := bufio.NewScanner(filePointer)

	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var isFirstLine bool = true
	var entries []Entry

	for scanner.Scan() {
		line := scanner.Bytes()

		// Check conditions
		if isFirstLine {
			isFirstLine = false
			continue
		}

		split_values := strings.Split(string(line), ",")

		if len(split_values) != 7 {
			fmt.Println("[ERROR]: Not enough values / Bad row")
			continue
		}

		entries = append(entries, Entry{
			Job_Number:       split_values[0],
			Job_Type:         split_values[1],
			Completed_By:     split_values[2],
			Start_Time:       parseTime(split_values[3]),
			Dispatch_Time:    parseTime(split_values[4]),
			In_Progress_Time: parseTime(split_values[5]),
			Complete_Time:    parseTime(split_values[6]),
		})
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("scanning file: %w", scanner.Err())
	}

	return entries, nil
}

func buildStartToDispatchWindows(entries []Entry) Windows {
	var windows Windows
	var nonDispatchedJobs Window

	for _, newEntry := range entries {
		var removeIndex []int
		for index, entry := range nonDispatchedJobs.Entries {
			if entry.Dispatch_Time.Before(newEntry.Start_Time) {
				removeIndex = append(removeIndex, index)
			}
		}

		if len(removeIndex) > 0 {
			nonDispatchedJobs.EndTime = newEntry.Start_Time
			windows.Windows = append(windows.Windows, Window{
				Entries:   slices.Clone(nonDispatchedJobs.Entries),
				StartTime: nonDispatchedJobs.StartTime,
				EndTime:   nonDispatchedJobs.EndTime,
			})

			slices.Reverse(removeIndex)
			for _, indexValue := range removeIndex {
				nonDispatchedJobs.Entries = slices.Delete(nonDispatchedJobs.Entries, indexValue, indexValue+1)
			}
		}

		if len(nonDispatchedJobs.Entries) == 0 {
			nonDispatchedJobs.StartTime = newEntry.Start_Time
		}

		nonDispatchedJobs.Entries = append(nonDispatchedJobs.Entries, newEntry)
	}

	if len(nonDispatchedJobs.Entries) > 0 {
		lastDispatch := nonDispatchedJobs.Entries[0].Dispatch_Time
		for _, entry := range nonDispatchedJobs.Entries[1:] {
			if entry.Dispatch_Time.After(lastDispatch) {
				lastDispatch = entry.Dispatch_Time
			}
		}

		nonDispatchedJobs.EndTime = lastDispatch
		windows.Windows = append(windows.Windows, Window{
			Entries:   slices.Clone(nonDispatchedJobs.Entries),
			StartTime: nonDispatchedJobs.StartTime,
			EndTime:   nonDispatchedJobs.EndTime,
		})
	}

	return windows
}

func StartToDispatch(fileName string) {
	entries, err := readEntries(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	windows := buildStartToDispatchWindows(entries)

	for _, value := range windows.Windows {
		fmt.Printf("%+v, ", value.StartTime)
		for _, tWindow := range value.Entries {
			fmt.Printf("%+v, ", tWindow.Job_Number)
		}
		fmt.Println("\n_______")
	}
}

func EveryMinute(fileName string) {

	file_pointer, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file_pointer.Close()

	scanner := bufio.NewScanner(file_pointer)

	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var windows Windows
	var currentTime time.Time
	var job Window

	one_minute, err := time.ParseDuration("60s")
	if err != nil {
		fmt.Println(err)
	}

	var isFirstLine bool = true
	var previous_start_time time.Time

	for scanner.Scan() {
		line := scanner.Bytes()

		split_values := strings.Split(string(line), ",")

		if isFirstLine {
			isFirstLine = false
			continue
		}

		if len(split_values) != 7 {
			fmt.Println("[ERROR]: Not enough values / Bad row")
			continue
		}

		var sheet_start_time = parseTime(split_values[3])

		if job.StartTime.IsZero() {
			currentTime = sheet_start_time
			job.StartTime = currentTime
		}

		if sheet_start_time.Sub(currentTime) > one_minute {
			job.EndTime = previous_start_time
			job.Job_Numbers = append(job.Job_Numbers, split_values[0])

			windows.Windows = append(windows.Windows, job)
			currentTime = sheet_start_time

			job = Window{}
			job.StartTime = sheet_start_time
			currentTime = sheet_start_time
			previous_start_time = sheet_start_time
		} else {
			job.Job_Numbers = append(job.Job_Numbers, split_values[0])
			previous_start_time = sheet_start_time
		}
	}

	windows.Windows = append(windows.Windows, job)

	if scanner.Err() != nil {
		fmt.Println("[ERROR]: scanning file", scanner.Err().Error())
		return
	}

	for _, value := range windows.Windows {
		fmt.Printf("%+v\n", value)
	}
}
