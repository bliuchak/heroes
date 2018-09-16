package heroes

import (
	"os"

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

// NewApplication returns pointer to App structure with filled data
func NewApplication(config config.Config) *App {
	return &App{
		Config: config,
	}
}

// InitLogger sets logger to App structure
func (a *App) InitLogger() {
	a.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}

// InitStorage sets database to App structure
func (a *App) InitStorage() error {
	s, err := db.NewRedis(a.Config.Database.Host, a.Config.Database.Password, a.Config.Database.Port)
	if err != nil {
		return err
	}
	a.Storage = s
	return nil
}

// Run runs server from App structure
func (a *App) Run() error {
	a.Logger.Info().Msg("Run app")

	a.Server = server.NewServer(a.Storage, a.Logger, a.Config)

	err := a.Server.Run()
	if err != nil {
		return err
	}
	return nil
}
