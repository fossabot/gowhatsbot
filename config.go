package main

import (
	"encoding/json"
	"io/fs"
	"os"
)

type Config struct {
	Config map[string]string
}

func newConfig(configpath string) Config {
	var c = Config{
		Config: map[string]string{},
	}
	c.LoadFromFile(configpath)
	return c
}

func (c *Config) SaveToFile(configpath string) error {
	if b, err := json.MarshalIndent(c.Config, "", "  "); err != nil {
		return err
	} else {
		return os.WriteFile(configpath, b, fs.ModeAppend)
	}
}

func (c *Config) LoadFromFile(configpath string) error {
	if cfg, err := os.ReadFile(configpath); err != nil {
		return err
	} else {
		return json.Unmarshal(cfg, &c.Config)
	}
}

func (c *Config) GetByKey(key, def string) string {
	if val, ok := c.Config[key]; ok {
		return val
	} else {
		c.Config[key] = def
		return def
	}
}

var WhatsConfig Config
var ConfigPath = "gowhatsbot.json"

func loadConfig() {

	WhatsConfig = newConfig(ConfigPath)
}
