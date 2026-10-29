package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var logRe = regexp.MustCompile(`(.{24}) \[(INFO|WARN|ERROR)] (.+): (.*)`)
var requestRe = regexp.MustCompile(`request_id=([a-zA-Z0-9_]+)`)
var userRe = regexp.MustCompile(`user_id=([a-zA-Z0-9_]+)`)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Service   string
	Message   string
	RequestID string
	UserID    string
}

func ParseLogLine(line string) (LogEntry, error) {
	line = strings.TrimSpace(line)

	matches := logRe.FindStringSubmatch(line)
	if matches == nil {
		return LogEntry{}, errors.New("log line is incorrect")
	}

	t, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		return LogEntry{}, err
	}
	message := matches[4]
	reqMatches := requestRe.FindStringSubmatch(message)
	if reqMatches == nil {
		return LogEntry{}, errors.New("log line is incorrect")
	}
	userMatches := userRe.FindStringSubmatch(message)
	userId := ""
	if len(userMatches) > 1 {
		userId = userMatches[1]
	}

	return LogEntry{t, matches[2], matches[3], message, reqMatches[1], userId}, nil
}

func ReadLogFile(filepath string) ([]LogEntry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	resLogEntries := make([]LogEntry, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lp, err := ParseLogLine(line)
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

func ScanLogDirectory(dirPath string) ([]string, error) {
	var logFiles []string

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".log") {
			logFiles = append(logFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return logFiles, nil
}

func ProcessMultipleFiles(filePaths []string) ([]LogEntry, error) {
	logEntries := make([]LogEntry, 0)
	for _, fp := range filePaths {
		le, err := ReadLogFile(fp)
		if err == nil {
			logEntries = append(logEntries, le...)
		}
	}

	return logEntries, nil
}

func CorrelateRequests(entries []LogEntry) map[string][]LogEntry {
	res := make(map[string][]LogEntry)
	for _, le := range entries {
		if le.RequestID == "" {
			res[""] = append(res[""], le)
		} else {
			res[le.RequestID] = append(res[le.RequestID], le)
		}
	}
	return res
}

func main() {
}
