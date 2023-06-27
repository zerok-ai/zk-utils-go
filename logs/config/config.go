package config

type LogsConfig struct {
	Color bool   `yaml:"color" env-description:"Add colors to the logs based on the log level"`
	Level string `yaml:"level" env-description:"Minimum log level to be allowed"`
}
