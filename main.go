package main

import (
	"fmt"
	"log-parser/internal/cli"
	"log-parser/internal/processor"
	"log-parser/internal/reporter"
	"log-parser/internal/scanner"
	"os"
	"time"
)

func main() {
	cfg, err := cli.ParseCommandLineArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	res := reporter.AnalysisResult{FailedRequests: make([]reporter.FailedRequestReport, 0)}
	start := time.Now()

	filePaths, err := scanner.ScanLogDirectory(cfg.InputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logEntries, err := processor.ProcessMultipleFiles(filePaths)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logsByReqId := processor.CorrelateRequests(logEntries)
	failedReqIds := processor.DetectFailedRequests(logsByReqId)

	for _, failedReqId := range failedReqIds {
		failedLog, _ := processor.FindFirstFailure(logsByReqId[failedReqId])
		timeline := make([]string, 0, len(logsByReqId[failedReqId]))
		for _, log := range processor.SortTimelineByTimestamp(logsByReqId[failedReqId]) {
			timeline = append(timeline, log.OriginalLog)
		}
		res.FailedRequests = append(res.FailedRequests, reporter.FailedRequestReport{
			RequestID:      failedLog.RequestID,
			FailingService: failedLog.Service,
			ErrorMessage:   failedLog.Message,
			Timeline:       timeline,
		})
	}

	res.TotalEntriesProcessed = len(logEntries)
	res.FailedRequestsFound = len(failedReqIds)
	res.ProcessingTimeSeconds = time.Since(start).Seconds()

	err = reporter.WriteJSONReport(res, cfg.OutputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
