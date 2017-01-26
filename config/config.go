// Package config provides configuration loading utilities.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the configurable service parameters.
type Config struct {
	RedisURL  string `toml:"redisURL"`
	SecretKey string `toml:"secretKey"`
}

// Load loads TOML file contents. If decoding fails, the service is aborted.
func (cfg *Config) Load(file string) {
	if _, err := toml.DecodeFile(file, &cfg); err != nil {
		fmt.Printf("Cannot load config due to %s", err)
		os.Exit(1)
	}
}
