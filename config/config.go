package config

import (
	"encoding/json"
	"github.com/tent/tent-client-go"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

var configPath string

func init() {
	configPath = os.Getenv("TENT_CONFIG")
	if configPath == "" {
		user, err := user.Current()
		if err != nil {
			fmt.Println(err)
			return
		}
		configPath = path.Join(user.HomeDir, ".config", "tent.json")
	}
}

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

func Read() (Config, error) {
	var config Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(file, &config); err != nil {
		return config, err
	}

	return config, nil
}