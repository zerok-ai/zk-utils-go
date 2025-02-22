package config

const LoggerTag = "badger_config"

type BadgerConfig struct {
	DBPath          string  `yaml:"badgerPath" env:"ZK_BADGER_PATH" env-description:"Badger DB Path"`
	BatchSize       int     `yaml:"batchSize"`
	GCDiscardRatio  float64 `yaml:"gcDiscardRatio"`
	GCTimerDuration int     `yaml:"gcTimerDuration"`
}
