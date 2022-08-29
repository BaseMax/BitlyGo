package main

import (
	"context"
	"log"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/itsjoniur/bitlygo/api"
	"github.com/itsjoniur/bitlygo/internal/configs"
	"github.com/itsjoniur/bitlygo/internal/durable"
)

var (
	configPath string = "config.yaml"
)

func main() {
	// Initialize configuration
	dir, err := configs.GetRootDir()
	if err != nil {
		logrus.Panicln(err)
	}

	configPath = path.Join(dir, configPath)
	if err := configs.Init(configPath); err != nil {
		log.Panicln(err)
	}

	configs := configs.AppConfig

	// Create a database client
	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     configs.Database.User,
		Password: configs.Database.Password,
		Host:     configs.Database.Host,
		Port:     configs.Database.Port,
		Name:     configs.Database.Name,
	})

	// Create logger
	logger := durable.NewLogger(logrus.New())

	// Serve HTTP
	if err := api.StartAPI(logger, db, configs.HTTP.Port); err != nil {
		log.Panicln(err)
	}
}
