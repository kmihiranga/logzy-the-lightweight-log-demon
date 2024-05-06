package log

import (
	"fmt"
	"sync"

	appConfig "logzy/config"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var logger *zap.Logger
var once sync.Once

func GetLogger() *zap.Logger {
	once.Do(func() {
		// read config content
		configFile := appConfig.CheckEnvironment("log")
		content := appConfig.ReadFile(configFile)
		appConfiguration := appConfig.Parse()

		var cfg zap.Config
		err := yaml.Unmarshal(*content, &cfg)
		if err != nil {
			panic(fmt.Sprintf("error unmarshal yaml: %v", err))
		}
		if len(cfg.OutputPaths) >= 2 {
			err = appConfig.CreateDirectoryIfNotExist(appConfiguration.AppConfig.LoggerFileLocation)
			if err != nil {
				panic(fmt.Sprintf("Error creating directory. %v", err))
			}
			fileCreated := appConfig.WriteFile(cfg.OutputPaths[1])
			if !fileCreated {
				panic("Log file not created for uber zap logger.")
			}
		}
		logger, err = cfg.Build()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()

		logger.Info("logger initialized")

	})
	return logger
}
