package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	ShopDB     string `yaml:"shop_db"`
	HardwareDB string `yaml:"hardware_db"`
	// DatabaseURL string `yaml:"database_url" env-required:"true"`

	KafkaBrokers []string `yaml:"kafka_brokers" env-required:"true"`
	KafkaTopic   string   `yaml:"kafka_topic" env-default:"product-updates"`

	GeneratorInterval  time.Duration `yaml:"generator_interval" env-default:"10s"`
	PriceChangePercent float64       `yaml:"price_change_percent" env-default:"10.0"`
	StockChangeAmount  int           `yaml:"stock_change_amount" env-default:"15"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
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
