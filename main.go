package main

import (
	"altimi-interview/file_logger"
	"altimi-interview/logger"
	"altimi-interview/sync"

	"flag"
	"fmt"
	"os"
)

func main() {
	sourceDir := flag.String("source", "", "Path to the source directory")
	targetDir := flag.String("target", "", "Path to the target directory")
	deleteMissing := flag.Bool("delete-missing", false, "Delete files from target that are missing in source")

	flag.Parse()

	if *sourceDir == "" || *targetDir == "" {
		fmt.Printf("Error: both --source and --target flags are required.")
		flag.Usage()
		os.Exit(1)
	}

	logger.UseLogger(file_logger.NewFileLogger("sync_log.txt"))
	defer logger.Close()

	if err := sync.SyncDirectories(*sourceDir, *targetDir, *deleteMissing); err != nil {
		logger.Message(fmt.Sprintf(logger.WARNING_LOG_COLOR+"Synchronization finished with errors: %v"+logger.RESET_LOG_COLOR, err))
	} else {
		logger.Message(logger.SUCCESS_LOG_COLOR + "Synchronization finished successfully." + logger.RESET_LOG_COLOR)
	}
}
