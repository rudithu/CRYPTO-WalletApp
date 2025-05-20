package config

import (
	"log"
	"os"
)

func InitLog() error {
	// Open a log file in append mode, create if not exists
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Print("ERROR: Failed to open log file: %v", err)
		return err
	}
	// Set output of logs to the file
	log.SetOutput(file)
	// Optional: log with date and time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return nil
}
