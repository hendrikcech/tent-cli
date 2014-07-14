package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/tent/hawk-go"
	"github.com/tent/tent-client-go"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
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

type SchemaConfig struct {
	Name     string `json:"name"`
	PostType string `json:"postType"`
}

// http://play.golang.org/p/__XjCi5hI7
func (s *SchemaConfig) MergeFragment(postType string) string {
	schema := strings.Split(s.PostType, "#")
	post := strings.Split(postType, "#")

	if len(post) > 1 {
		return schema[0] + "#" + post[1]
	}
	return s.PostType
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
	Schemas  []SchemaConfig  `json:"schemas"`
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
		c.setDefault()
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

func (c *Config) setDefault() {
	c.Schemas = []SchemaConfig{
		SchemaConfig{"meta", "https://tent.io/types/meta/v0"},
		SchemaConfig{"app", "https://tent.io/types/app/v0"},
		SchemaConfig{"status", "https://tent.io/types/status/v0"},
	}
}

func (c *Config) ProfileByName(name string) (int, *ProfileConfig) {
	for i, p := range c.Profiles {
		if p.Name == name {
			return i, &c.Profiles[i]
		}
	}
	return -1, &ProfileConfig{}
}

func (c *Config) DefaultProfile() (*ProfileConfig, error) {
	if c.Default == "" {
		return &ProfileConfig{}, fmt.Errorf("No default profile set.")
	}

	i, p := c.ProfileByName(c.Default)
	if i == -1 {
		return &ProfileConfig{}, fmt.Errorf("Default profile \"%v\" doesn't exist.", c.Default)
	}

	return p, nil
}

func (c *Config) SchemaByName(name string) (int, *SchemaConfig) {
	n := strings.Split(name, "#")
	for i, s := range c.Schemas {
		if s.Name == n[0] {
			return i, &c.Schemas[i]
		}
	}
	return -1, &SchemaConfig{}
}
