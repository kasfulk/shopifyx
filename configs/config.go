package configs

import (
	"os"
)

type Server struct {
	Host    string
	Port    string
	Version string
}

type Database struct {
	Host string
	Port string
	Name string
	User string
	Pass string
}

type Config struct {
	Server   *Server
	Database *Database
}

func LoadConfig() *Config {
	var (
		dbuser     = "postgres"
		dbpassword = "postgres"
		dbport     = "5432"
		dbhost     = "localhost"
		dbname     = "shopifyx"
	)

	var (
		host    = "0.0.0.0"
		port    = "5000"
		version = "v1"
	)

	if os.Getenv("SYSTEMD_DB_NAME") != "" {
		dbuser = os.Getenv("SYSTEMD_DB_USER")
		dbpassword = os.Getenv("SYSTEMD_DB_PASSWORD")
		dbhost = os.Getenv("SYSTEMD_DB_HOST")
		dbport = os.Getenv("SYSTEMD_DB_PORT")
		dbname = os.Getenv("SYSTEMD_DB_NAME")
	}

	if os.Getenv("SYSTEMD_APP_ENV") != "" {
		host = os.Getenv("SYSTEMD_APP_HOST")
		port = os.Getenv("SYSTEMD_APP_PORT")
		version = os.Getenv("SYSTEM_APP_VERSION")
	}

	var (
		config = &Config{
			Database: &Database{
				Host: dbhost,
				Port: dbport,
				Name: dbname,
				User: dbuser,
				Pass: dbpassword,
			},
			Server: &Server{
				Host:    host,
				Port:    port,
				Version: version,
			},
		}
	)

	return config
}
