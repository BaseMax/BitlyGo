package configs

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Struct for configs
type Config struct {
	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http"`
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	}
}

var AppConfig *Config

// Init initialize project configuration
func Init(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var configs map[string]Config
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return err
	}

	cfg := configs["default"]
	AppConfig = &cfg
	return nil
}

// GetRootDir find the root directory of project
func GetRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}
