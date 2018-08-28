package heroes

import (
	"os"

	"github.com/bliuchak/heroes/internal/config"
	"github.com/bliuchak/heroes/internal/db"
	"github.com/bliuchak/heroes/internal/server"
	"github.com/bliuchak/heroes/internal/storage"
	"github.com/rs/zerolog"
)

// App ...
type App struct {
	Logger zerolog.Logger
	Storage storage.Storager
	Server *server.Server
	Config config.Config
}

// NewApplication ...
func NewApplication(config config.Config) *App {
	return &App{
		Config: config,
	}
}

func (a *App) InitLogger() {
	a.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}

func (a *App) InitStorage() error {
	s, err := db.NewRedis(a.Config.Database.Host, a.Config.Database.Password, a.Config.Database.Port)
	if err != nil {
		return err
	}
	a.Storage = s
	return nil
}

func (a *App) Run() error {
	a.Logger.Info().Msg("Run app")

	a.Server = server.NewServer(a.Storage, a.Logger, a.Config)

	err := a.Server.Run()
	if err != nil {
		return err
	}
	return nil
}
