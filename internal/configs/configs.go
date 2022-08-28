package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
