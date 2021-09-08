package config

import (
	"Oracle-Hackathon-BE/database"
	"log"

	"github.com/spf13/viper"
)

var (
	CFG = &Config{}
)

type Config struct {
	reader *viper.Viper
}

func New() *Config {

	viper := viper.New()
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	CFG = &Config{reader: viper}
	return CFG
}

func (c *Config) ReadEnv(key string) string {
	return CFG.reader.GetString(key)
}

func (c *Config) FetchDatabaseConfig() (db *database.DatabaseConfig) {
	return &database.DatabaseConfig{}
}
