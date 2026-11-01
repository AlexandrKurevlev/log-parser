package reporter

import (
	"encoding/json"
	"os"
)

type AnalysisResult struct {
	TotalEntriesProcessed int                   `json:"total_entries_processed"`
	FailedRequestsFound   int                   `json:"failed_requests_found"`
	ProcessingTimeSeconds float64               `json:"processing_time_seconds"`
	FailedRequests        []FailedRequestReport `json:"failed_requests"`
}

type FailedRequestReport struct {
	RequestID      string   `json:"request_id"`
	FailingService string   `json:"failing_service"`
	ErrorMessage   string   `json:"error_message"`
	Timeline       []string `json:"timeline"`
}

func WriteJSONReport(result AnalysisResult, filename string) error {
	data, _ := json.MarshalIndent(result, "", "  ")
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
