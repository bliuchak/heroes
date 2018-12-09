package config

// Config contains application config data
type Config struct {
	Database Database
	Server   Server
}

// Database contains database config data
type Database struct {
	Host     string
	Port     string
	Password string
}

// Server contains server config data
type Server struct {
	Port int
}

// NewConfig returns pointer on Config with filled data
func NewConfig(appport int, dbhost string, dbport string, dbpassword string) *Config {
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
