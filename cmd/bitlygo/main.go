package main

import (
	"context"
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/itsjoniur/bitlygo/api"
	"github.com/itsjoniur/bitlygo/internal/configs"
	"github.com/itsjoniur/bitlygo/internal/durable"
	"github.com/sirupsen/logrus"
)

var (
	configPath string = "../internal/configs/config.yaml"
)

func main() {
	// initialize configuration
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(path.Join(path.Dir(b)))

	configPath = path.Join(dir, configPath)
	if err := configs.Init(configPath); err != nil {
		log.Panicln(err)
	}

	configs := configs.AppConfig
	// create a database client
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     configs.Database.User,
		Password: configs.Database.Password,
		Host:     configs.Database.Host,
		Port:     configs.Database.Port,
		Name:     configs.Database.Name,
	})
	// create logger
	logger := durable.NewLogger(logrus.New())
	// serve HTTP
	if err := api.StartAPI(logger, db, configs.HTTP.Port); err != nil {
		log.Panicln(err)
	}
}
