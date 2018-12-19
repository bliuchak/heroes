package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bliuchak/heroes/internal"
	"github.com/bliuchak/heroes/internal/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	appport    = kingpin.Flag("appport", "port where to run app").Envar("APP_PORT").Default("3001").Int()
	dbhost     = kingpin.Flag("dbhost", "storage host").Envar("DB_HOST").String()
	dbport     = kingpin.Flag("dbport", "storage port").Envar("DB_PORT").String()
	dbpassword = kingpin.Flag("dbpassword", "storage password").Envar("DB_PASSWORD").String()
)

func main() {
	kingpin.Parse()

	conf := config.NewConfig(*appport, *dbhost, *dbport, *dbpassword)
	app := heroes.NewApplication(*conf)

	app.InitLogger()
	err := app.InitStorage()
	if err != nil {
		app.Logger.Error().Err(err).Msg("Unable to init storage")
	}

	// create channel to collect terminate or interrupt signals
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		timeout := 5 * time.Second
		app.Logger.Info().Interface("sig", sig).Dur("timeout", timeout).Msg("Graceful shutdown with timeout")
		err := app.Server.Stop(timeout)
		if err != nil {
			app.Logger.Error().Err(err).Msg("Graceful shutdown error")
		}
		app.Logger.Info().Interface("sig", sig).Msg("Server stopped")
	}()

	err = app.Run()
	if err != nil {
		app.Logger.Error().Err(err).Msg("Unable to run app")
	}
}
