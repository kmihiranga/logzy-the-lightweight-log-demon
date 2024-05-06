package main

import (
	"fmt"
	"sync"

	"logzy/config"
	logger "logzy/pkg/log"
	tailLog "logzy/pkg/tail"

	"go.uber.org/zap"
)

func main() {
	var wg sync.WaitGroup
	var log *zap.SugaredLogger = logger.GetLogger().Sugar()

	cfg := config.Parse()

	logFileName := fmt.Sprintf(cfg.AppConfig.LogLocation + "/" + cfg.AppConfig.LogFileName)
	if logFileName == "" {
		log.Error("cannot find the log file to watch")
		panic(fmt.Errorf("cannot find the log file to watch"))
	}
	filesToTail := []string{logFileName} // List of log files

	for _, filePath := range filesToTail {
		offsetFilePath := filePath + ".offset" // Separate offset file for each log
		fileCreated := config.WriteFile(offsetFilePath)
		if !fileCreated {
			panic("offset file not created")
		}
		wg.Add(1)
		go tailLog.TailLogFile(filePath, offsetFilePath, &wg)
	}

	wg.Wait() // Wait for all tailing goroutines to complete
}
