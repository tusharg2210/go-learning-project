package config

import (
	"flag"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	address string
}

type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true" envDefault:"development"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true" envDefault:"./storage"`
	HTTP HTTPConfig `yaml:"http_server"`

}

func MustLoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if(configPath == "") {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path isnot set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist at path: %s", configPath)
	}

	var cfg Config
	
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	return &cfg
}	