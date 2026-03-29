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

func StartToDispatch(fileName string) {

	filePointer, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Unable to open file", err)
		return
	}

	defer filePointer.Close()

	scanner := bufio.NewScanner(filePointer)

	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var windows Windows
	// var currentTime time.Time
	// var previous_start_time time.Time
	var non_dispatched_jobs Window
	var isFirstLine bool = true

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

		newEntry := Entry{
			Job_Number:       split_values[0],
			Job_Type:         split_values[1],
			Completed_By:     split_values[2],
			Start_Time:       parseTime(split_values[3]),
			Dispatch_Time:    parseTime(split_values[4]),
			In_Progress_Time: parseTime(split_values[5]),
			Complete_Time:    parseTime(split_values[6]),
		}

		// Main logic
		var allReadyAdded bool
		var removeIndex []int
		for index, entry := range non_dispatched_jobs.Entries {
			// IF current window has a job/s that got dispatched
			if entry.Dispatch_Time.Sub(newEntry.Start_Time) < 0 {

				// Add endtime to be dispatch time
				non_dispatched_jobs.EndTime = newEntry.Start_Time

				if !allReadyAdded {
					// Append window to windows
					entriesSnapshot := slices.Clone(non_dispatched_jobs.Entries)
					windows.Windows = append(windows.Windows, Window{
						Entries:   entriesSnapshot,
						StartTime: non_dispatched_jobs.StartTime,
						EndTime:   non_dispatched_jobs.EndTime,
					})
					allReadyAdded = true

					non_dispatched_jobs.StartTime = non_dispatched_jobs.Entries[0].Start_Time
				}

				removeIndex = append(removeIndex, index)
			}
		}

		// Reverse the list of index-es so it deletes it backwards so the when trying to index in to remove and doesn't affect the current index count
		slices.Reverse(removeIndex)
		for _, index_value := range removeIndex {
			non_dispatched_jobs.Entries = slices.Delete(non_dispatched_jobs.Entries, index_value, index_value+1)
		}

		// Append current job to window
		non_dispatched_jobs.Entries = append(non_dispatched_jobs.Entries, newEntry)

	}

	for _, value := range windows.Windows {
		fmt.Printf("%+v, ", value.StartTime)
		for _, t_window := range value.Entries {
			fmt.Printf("%+v, ", t_window.Job_Number)
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
