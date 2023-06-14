package config

type RedisConfig struct {
	Host        string         `yaml:"host" env:"REDIS_HOST" env-description:"Database host"`
	Port        string         `yaml:"port" env:"REDIS_PORT" env-description:"Database port"`
	DBs         map[string]int `yaml:"dbs" env:"REDIS_DB" env-description:"Database to load"`
	ReadTimeout int            `yaml:"readTimeout"`
}

type DB struct {
	Name string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
}
