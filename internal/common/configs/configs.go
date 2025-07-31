package configs

import "time"

type Config struct {
	Storage     Storage     `envconfig:"STORAGE" required:"true"`
	Server      Server      `envconfig:"SERVER" required:"true"`
	Logger      Logger      `envconfig:"LOGGER" required:"true"`
	Middlerware Middlerware `envconfig:"MIDDLEWARE" required:"true"`
}

func init() {
	mustLoadCfg(singleConfig)
}

var singleConfig = &Config{}

func Get() *Config {
	return singleConfig
}

type Storage struct {
	Host     string `envconfig:"HOST"`
	Port     string `envconfig:"PORT" default:"5432"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASS"`
	Name     string `envconfig:"NAME"`
	SSLmode  string `envconfig:"SSLM" default:"disable"`
}

type Server struct {
	Address           string        `envconfig:"ADDR"`
	Port              string        `envconfig:"PORT"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT" default:"5s"`
}

type Logger struct {
	Level      string `envconfig:"LEVEL"`
	Encoding   string `envconfig:"ENCODING" default:"json"`
	Output     string `envconfig:"OUTPUT" default:"stdout"`
	MessageKey string `envconfig:"MESSAGE_KEY" default:"message"`
}

type Middlerware struct {
	AuthToken string `envconfig:"AUTH_TOKEN"`
}
