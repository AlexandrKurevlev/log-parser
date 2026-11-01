package processor

import (
	"bufio"
	"fmt"
	"log-parser/internal/parser"
	"os"
	"sort"
)

func ReadLogFile(filepath string) ([]parser.LogEntry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resLogEntries := make([]parser.LogEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lp, err := parser.ParseLogLine(line)
		if err != nil {
			fmt.Println(err)
		} else {
			resLogEntries = append(resLogEntries, lp)
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return resLogEntries, nil
}

func ProcessMultipleFiles(filePaths []string) ([]parser.LogEntry, error) {
	logEntries := make([]parser.LogEntry, 0)
	for _, fp := range filePaths {
		le, err := ReadLogFile(fp)
		if err == nil {
			logEntries = append(logEntries, le...)
		}
	}

	return logEntries, nil
}

func CorrelateRequests(entries []parser.LogEntry) map[string][]parser.LogEntry {
	res := make(map[string][]parser.LogEntry)
	for _, le := range entries {
		if le.RequestID == "" {
			res[""] = append(res[""], le)
		} else {
			res[le.RequestID] = append(res[le.RequestID], le)
		}
	}
	return res
}

func DetectFailedRequests(correlatedRequests map[string][]parser.LogEntry) []string {
	res := make([]string, 0)
	for reqId, entries := range correlatedRequests {
		for _, entry := range entries {
			if entry.Level == "WARN" || entry.Level == "ERROR" {
				res = append(res, reqId)
				break
			}
		}
	}
	return res
}

func FindFirstFailure(requestEntries []parser.LogEntry) (parser.LogEntry, bool) {
	sort.Slice(requestEntries, func(i, j int) bool {
		return requestEntries[i].Timestamp.Before(requestEntries[j].Timestamp)
	})

	for _, entry := range requestEntries {
		if entry.Level == "WARN" || entry.Level == "ERROR" {
			return entry, true
		}
	}

	return parser.LogEntry{}, false
}

func SortTimelineByTimestamp(entries []parser.LogEntry) []parser.LogEntry {
	cloned := make([]parser.LogEntry, len(entries))
	copy(cloned, entries)
	sort.Slice(cloned, func(i, j int) bool {
		return cloned[i].Timestamp.Before(cloned[j].Timestamp)
	})
	return cloned
}
