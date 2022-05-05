package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type db struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DBName   string `toml:"dbname"`
	User     string `toml:"user"`
	Password string `toml:"pass"`
}

type main struct {
	Providers []string `toml:"providers"`
	Fsyms     []string `toml:"fsyms"`
	Tsyms     []string `toml:"tsyms"`
}

type Config struct {
	DB   db   `toml:"db"`
	Main main `toml:"main"`
}

func NewAppConfig(filename string) (*Config, error) {
	var config *Config
	if filename == "" {
		return &Config{
			DB: db{
				Host:     "localhost",
				Port:     5432,
				DBName:   "ticker",
				User:     "ticker",
				Password: "ticker",
			},
			Main: main{
				Providers: []string{"cryptocompare"},
				Fsyms:     []string{"BTC", "XRP"},
				Tsyms:     []string{"USD", "EUR"},
			},
		}, nil
	} else {
		if _, err := toml.DecodeFile(filename, &config); err != nil {
			return nil, fmt.Errorf("failed to decode file %s, %s", filename, err)
		}
	}

	return config, nil
}
