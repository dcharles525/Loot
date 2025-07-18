package main

import (
	"fmt"
	"log"
	"os"

	"loot/ui"
)

func main() {
	// Setup logging early
	logFile, err := os.OpenFile("loot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if err := ui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "loot: %v\n", err)
		os.Exit(1)
	}
}
