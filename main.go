package main

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var logRe = regexp.MustCompile(`(.{24}) \[(INFO|WARN|ERROR)] (.+): (.*)`)
var requestRe = regexp.MustCompile(`request_id=([a-zA-Z0-9_]+)`)
var userRe = regexp.MustCompile(`user_id=([a-zA-Z0-9_]+)`)

type LogEntry struct {
	TimeStamp time.Time
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

func main() {
}
