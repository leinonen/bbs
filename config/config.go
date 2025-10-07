package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ListenAddr     string `json:"listen_addr"`
	DatabasePath   string `json:"database_path"`
	ServerName     string `json:"server_name"`
	HostKeyPath    string `json:"host_key_path"`
	AllowAnonymous bool   `json:"allow_anonymous"`
	MaxUsers       int    `json:"max_users"`
}

func Default() *Config {
	return &Config{
		ListenAddr:     ":2222",
		DatabasePath:   "bbs.db",
		ServerName:     "Go BBS System",
		HostKeyPath:    "host_key",
		AllowAnonymous: true,
		MaxUsers:       100,
	}
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}