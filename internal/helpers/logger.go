package helpers

import (
	"log"
	"os"
	"path/filepath"
)

func SetupLogging() *os.File {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal("Error creating log directory:", err)
	}

	logFile, err := os.OpenFile(filepath.Join("logs", "server.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	log.SetOutput(logFile)

	// log.Lshortfile())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile
}
