package config

type Config struct {
	AppConfig *AppConfig
}

type AppConfig struct {
	Name               string   `yaml:"name"`
	Env                string   `yaml:"env"`
	SlackUri           string   `yaml:"slack-uri"`
	TargetService      string   `yaml:"target-service"`
	LogFileName        string   `yaml:"log-file-name"`
	LogLocation        string   `yaml:"log-location"`
	LoggerFileLocation string   `yaml:"logger-file-location"`
	ResponseTimeout    int      `yaml:"response-timeout"`
	StartErrorPatterns []string `yaml:"start-error-patterns"`
}
