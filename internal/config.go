package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Timeout      float32 `yaml:"timeout"`
	RetriesCount uint16  `yaml:"retriesCount"`
	WorkersCount uint16  `yaml:"workersCount"`
}

var (
	ConfigPathEnv = "CONFIG_PATH"
)

func MustLoad() *Config {
	configPath, _ := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() (string, error) {
	const configFile = ".env"

	err := godotenv.Load(configFile)
	if err != nil {
		return "", err
	}

	cfgPath, exists := os.LookupEnv(ConfigPathEnv)
	if !exists {
		return "", fmt.Errorf("config wasn't found")
	}

	return cfgPath, nil
}
