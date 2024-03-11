package configs

import (
	"os"
	"strconv"
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

type Auth struct {
	Secret string
	Salt   int
}

type Config struct {
	Server   *Server
	Database *Database
	Auth     *Auth
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
		port    = "8000"
		version = "v1"
	)

	var (
		secret = "secret"
		salt   = int(8)
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

	if os.Getenv("SYSTEMD_APP_AUTH_SECRET") != "" {
		secret = os.Getenv("SYSTEMD_APP_AUTH_SECRET")
		salt, _ = strconv.Atoi(os.Getenv("SYSTEMD_APP_AUTH_SALT"))
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
			Auth: &Auth{
				Secret: secret,
				Salt:   salt,
			},
		}
	)

	return config
}
