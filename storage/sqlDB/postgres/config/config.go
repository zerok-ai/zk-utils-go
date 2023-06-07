package config

type PostgresConfig struct {
	Host                           string `yaml:"host" env-description:"Postgres host"`
	Port                           int    `yaml:"port" env-description:"Postgres port"`
	User                           string `yaml:"user" env:"PL_POSTGRES_USERNAME" env-description:"Postgres user"`
	Password                       string `yaml:"password" env:"PL_POSTGRES_PASSWORD" env-description:"Postgres password"`
	Dbname                         string `yaml:"dbname" env-description:"Postgres database"`
	MaxConnections                 int    `json:"max_connections"`
	MaxIdleConnections             int    `json:"max_idle_connections"`
	ConnectionMaxLifetimeInMinutes int    `json:"connection_max_lifetime_in_minutes"`
}
