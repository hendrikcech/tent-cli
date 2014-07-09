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
	"path"
	"runtime"
)

func configPath() string {
	p := os.Getenv("TENT_CONFIG")
	if p == "" {
		p = path.Join(homedir(), ".config", "tent.json")
	}
	return p
}

func homedir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("%APPDATA%")
	}
	return os.Getenv("HOME")
}

type ProfileConfig struct {
	Name    string                `json:"name"`
	Entity  string                `json:"entity"`
	Servers []tent.MetaPostServer `json:"servers"`

	ID  string `json:"id"`
	Key string `json:"key"`
	App string `json:"app"`
}

func (p *ProfileConfig) Client() *tent.Client {
	var c *hawk.Credentials
	if p.ID != "" {
		c = &hawk.Credentials{
			ID:   p.ID,
			Key:  p.Key,
			App:  p.App,
			Hash: sha256.New,
		}
	}
	return &tent.Client{
		Credentials: c,
		Servers:     p.Servers,
		Entity:      p.Entity,
	}
}

type Config struct {
	Default  string          `json:"default"`
	Profiles []ProfileConfig `json:"profiles"`
}

func (c *Config) Write() error {
	enc, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	p := configPath()

	err = os.MkdirAll(path.Dir(p), 0700)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p, enc, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Read() error {
	p := configPath()

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil
	}

	file, err := ioutil.ReadFile(p)
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
		return &ProfileConfig{}, errors.New("No default profile set.")
	}

	i, p := c.ByName(c.Default)
	if i == -1 {
		err := errors.New(fmt.Sprintf("Default profile \"%v\" doesn't exist.", c.Default))
		return &ProfileConfig{}, err
	}

	return p, nil
}
