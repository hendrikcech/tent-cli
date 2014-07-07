package config

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tent/hawk-go"
	"github.com/tent/tent-client-go"
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

type ProfileConfig struct {
	Name    string
	Entity  string
	Servers []tent.MetaPostServer

	ID  string
	Key string
	App string
}

func (p *ProfileConfig) Client() *tent.Client {
	return &tent.Client{
		Credentials: &hawk.Credentials{
			ID:   p.ID,
			Key:  p.Key,
			App:  p.App,
			Hash: sha256.New,
		},
		Servers: p.Servers,
	}
}

type Config struct {
	Default  string
	Profiles []ProfileConfig
}

func (c *Config) Write() error {
	enc, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(configPath), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, enc, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Read() error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(file, &c); err != nil {
		return err
	}

	return nil
}

func (c *Config) ByName(name string) (int, *ProfileConfig) {
	for i, p := range c.Profiles {
		if p.Name == name {
			return i, &c.Profiles[i]
		}
	}
	return -1, &ProfileConfig{}
}

func (c *Config) DefaultProfile() (*ProfileConfig, error) {
	if c.Default == "" {
		return &ProfileConfig{}, errors.New("no default profile set")
	}

	i, p := c.ByName(c.Default)
	if i == -1 {
		err := errors.New(fmt.Sprintf("default profile \"%v\" doesn't exist", c.Default))
		return &ProfileConfig{}, err
	}

	return p, nil
}
