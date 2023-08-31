package config

type RedisConfig struct {
	Host        string         `yaml:"host"`
	Port        string         `yaml:"port"`
	DBs         map[string]int `yaml:"dbs"`
	ReadTimeout int            `yaml:"readTimeout"`
	Password    string         `yaml:"password" env:"REDIS_PASSWORD" env-description:"Redis password"`
}

type DB struct {
	Name string `yaml:"host" env:"DB_HOST" env-description:"Database host"`
}
