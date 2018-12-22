package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bliuchak/heroes/internal"
	"github.com/bliuchak/heroes/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	appport    = kingpin.Flag("appport", "port where to run app").Envar("APP_PORT").Default("3001").Int()
	dbhost     = kingpin.Flag("dbhost", "storage host").Envar("DB_HOST").String()
	dbport     = kingpin.Flag("dbport", "storage port").Envar("DB_PORT").String()
	dbpassword = kingpin.Flag("dbpassword", "storage password").Envar("DB_PASSWORD").String()
)

func init() {
	kingpin.Parse()
}

func main() {
	conf := config.NewConfig(*appport, *dbhost, *dbport, *dbpassword)
	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	app := heroes.NewApplication(*conf, logger)

	err := app.InitStorage()
	if err != nil {
		app.Logger.Error().Err(err).Msg("Unable to init storage")
	}

	app.InitServer()

	// create channel to collect terminate or interrupt signals
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		app.Logger.Info().Int("port", app.Config.Server.Port).Msg("Start app")
		serverErr <- app.Server.Start()
	}()

	select {
	case err = <-serverErr:
		close(serverErr)
	case sig := <-gracefulShutdown:
		close(gracefulShutdown)
		timeout := 5 * time.Second
		app.Logger.Warn().Interface("sig", sig).Dur("timeout", timeout).Msg("Start graceful shutdown with timeout")
		err = app.Server.Stop(timeout)
		if err == nil {
			app.Logger.Warn().Msg("Finish graceful shutdown")
		}
	}

	if err != nil {
		app.Logger.Error().Err(err).Msg("Server stopped with error")
	}
}
