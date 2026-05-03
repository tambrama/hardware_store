package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	DatabaseURL string `yaml:"database_url" env-required:"true"`

	HTTPServer `yaml:"http_server"`
	Clients    ClientsConfig `yaml:"clients"`

	AppSecret string `yaml:"jwt_secret_key" env:"APP_SECRET"`
	AppID     string `yaml:"app_id"`

	KafkaBrokers []string `yaml:"kafka_brokers" env-required:"true"`
	KafkaTopic   string   `yaml:"kafka_topic" env-default:"product-updates"`

	RedisHost string `yaml:"redis_host" env-default:"redis"`
	RedisPort string `yaml:"redis_port" env-default:"6379"`

	PhotoServiceURL string `env:"PHOTO_SERVICE_URL,required"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriers_count"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set ")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	return &cfg
}
