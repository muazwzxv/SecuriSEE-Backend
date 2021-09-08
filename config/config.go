package config

import (
	"log"

	"github.com/spf13/viper"
)

var CFG Config

type Config struct {
	reader *viper.Viper
}

func New() Config {

	if CFG.reader == nil {

		viper := viper.New()
		viper.AddConfigPath("./")
		viper.SetConfigFile("config.yml")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error while reading config file %s", err)
		}

		CFG = Config{reader: viper}
	}

	return CFG
}

// Return global instance
func (c *Config) GetInstance() Config {
	return CFG
}

func (c *Config) readEnv(key string) string {
	return CFG.reader.GetString(key)
}

func (c *Config) FetchDatabaseConfig() *DatabaseConfig {

	port := CFG.reader.GetInt("Database.Port")

	return &DatabaseConfig{
		User:         CFG.readEnv("Database.User"),
		Password:     CFG.readEnv("Database.Password"),
		Host:         CFG.readEnv("Database.Host"),
		DatabaseName: CFG.readEnv("Database.Name"),
		Port:         port,
	}

}
