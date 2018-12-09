package main

import (
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

	err = app.Run()
	if err != nil {
		app.Logger.Error().Err(err).Msg("Unable to run app")
	}
}
