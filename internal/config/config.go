package config

type Config struct {
	Database Database
	Server   Server
}

type Database struct {
	Host     string
	Port     int
	Password string
}

type Server struct {
	Port int
}

func NewConfig(appport int, dbhost string, dbport int, dbpassword string) *Config {
	return &Config{
		Database: Database{
			Host:     dbhost,
			Port:     dbport,
			Password: dbpassword,
		},
		Server: Server{
			Port: appport,
		},
	}
}
