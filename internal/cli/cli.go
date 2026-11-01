package cli

import (
	"flag"
	"os"
)

type Config struct {
	InputDir   string
	OutputFile string
}

func ParseCommandLineArgs() (Config, error) {
	inputDir := flag.String("input-dir", ".", "Directory containing .log files")
	outputFile := flag.String("output-file", "results.json", "JSON output file path")
	flag.Parse()

	_, err := os.Stat(*inputDir)
	if err != nil {
		return Config{}, err
	}

	return Config{*inputDir, *outputFile}, nil
}
