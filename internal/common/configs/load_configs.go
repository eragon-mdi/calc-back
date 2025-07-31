package configs

import (
	"log"

	"github.com/go-faster/errors"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	EnvPath                  = ".env"
	ErrLoadCfgFile           = "failed to load dot-env file"
	ErrUnmarshalCfgsFromFile = "failed to load configurations"
)

func mustLoadCfg(c *Config) {
	var errLoad error
	if err := godotenv.Load(EnvPath); err != nil {
		errLoad = errors.Wrap(err, ErrLoadCfgFile) // in case with docker-compose, no need to load
	}

	if err := envconfig.Process("", c); err != nil {
		log.Fatal(errors.Join(errors.Wrap(err, ErrUnmarshalCfgsFromFile), errLoad))
	}
}
