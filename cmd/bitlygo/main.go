package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"

	"github.com/itsjoniur/bitlygo/internal/configs"
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
	// arguments: file, environment
	// create a database client
	// arguments: connection info
	// create logger
	// serve HTTP
	err := http.ListenAndServe(fmt.Sprintf(":%s", configs.HTTP.Port), nil)
	if err != nil {
		log.Panicln(err)
	}
}
