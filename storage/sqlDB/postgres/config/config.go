package config

type PostgresConfig struct {
	Host     string `yaml:"host" env-description:"Postgres host"`
	Port     int    `yaml:"port" env-description:"Postgres port"`
	User     string `yaml:"user" env:"PL_POSTGRES_USERNAME" env-description:"Postgres user"`
	Password string `yaml:"password" env:"PL_POSTGRES_PASSWORD" env-description:"Postgres password"`
	Dbname   string `yaml:"dbname" env-description:"Postgres database"`
}
