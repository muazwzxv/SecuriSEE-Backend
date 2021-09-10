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

	port := CFG.reader.GetInt("Development.Database.Port")

	return &DatabaseConfig{
		User:         CFG.readEnv("Development.Database.User"),
		Password:     CFG.readEnv("Development.Database.Password"),
		Host:         CFG.readEnv("Development.Database.Host"),
		DatabaseName: CFG.readEnv("Development.Database.Name"),
		Port:         port,
	}

}
