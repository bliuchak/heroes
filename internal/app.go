package heroes

import (
	"github.com/bliuchak/heroes/internal/config"
	"github.com/bliuchak/heroes/internal/db"
	"github.com/bliuchak/heroes/internal/server"
	"github.com/bliuchak/heroes/internal/storage"
	"github.com/rs/zerolog"
)

// App is an application container with necessary dependencies
type App struct {
	Logger  zerolog.Logger
	Storage storage.Storager
	Server  server.Serverer
	Config  config.Config
}

// NewApplication returns pointer to App structure with config and logger
func NewApplication(c config.Config, l zerolog.Logger) *App {
	return &App{
		Config: c,
		Logger: l,
	}
}

// InitStorage initialize application storage
func (a *App) InitStorage() error {
	s, err := db.NewRedis(a.Config.Database.Host, a.Config.Database.Password, a.Config.Database.Port)
	if err != nil {
		return err
	}
	a.Storage = s
	return nil
}

// InitServer initialise application server
func (a *App) InitServer() {
	a.Server = server.NewServer(a.Storage, a.Logger, a.Config)
}
