package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"loot/ui"
)

func main() {
	//Ensure the .loot directory is setup
	homeDir, _ := os.UserHomeDir()
	mkdirErr := os.MkdirAll(filepath.Join(homeDir, ".loot"), 755)
	if mkdirErr != nil {
		log.Fatalf("Error creating directory: %v", mkdirErr)
	}

	//Create log file
	logFile, logFileErr := os.OpenFile(
		filepath.Join(homeDir, ".loot", "loot.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if logFileErr != nil {
		log.Fatalf("failed to open log file: %v", logFileErr)
	}
	defer logFile.Close()

	//Route all output to log file
	log.SetOutput(logFile)

	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "loot: %v\n", err)
		os.Exit(1)
	}
}
