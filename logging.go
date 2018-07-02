package main

import (
	"io"
	"log"
	"os"
)

func SetupLog(component string, logLevel string, logFile string, logStderr bool) {
	//	log.SetSource(component)
	//	log.Infof("Set log level to %s", logLevel)
	//	log.SetLevel(logLevel)

	outputWriters := make([]io.Writer, 0, 2)
	if logFile != "" {
		fileWriter, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("open log file %s failed: %v", logFile, err)
		}
		outputWriters = append(outputWriters, fileWriter)
	}

	if logStderr {
		outputWriters = append(outputWriters, os.Stderr)
	}

	if len(outputWriters) == 1 {
		log.SetOutput(outputWriters[0])
	} else {
		log.SetOutput(io.MultiWriter(outputWriters...))
	}
}
