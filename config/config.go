package config

import (
	"encoding/json"
	"github.com/tent/tent-client-go"
	// "fmt"
)

type EntityConfig struct {
	Entity string
	Servers []tent.MetaPostServer

	ID string
	Key string
	App string

	Short string

}


type Config struct {
	Entities []EntityConfig
}

func Write() error {
	return nil
}

var data = ``

func Read() (Config, error) {
	var config Config

	err := json.Unmarshal([]byte(data), &config)

	if err != nil {
		return config, err
	}

	return config, nil
}