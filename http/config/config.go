package config

var HttpDebug = false

func Init(httpDebug bool) {
	HttpDebug = httpDebug
}

type HttpConfig struct {
	Debug bool `yaml:"debug" env:"HTTP_DEBUG" env-description:"Whether to pass debug information in the response or not"`
}
