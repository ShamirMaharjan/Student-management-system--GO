package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-reqired:"true"`
	StoragePath string `yaml:"storage_path" env-reqired:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("configuration path is not set")
		}
	}

	//check if the path exists
	//no use of info so we use _
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Congiguration file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("cannot read configuration file: %s", err.Error())
	}

	return &cfg
}
