package config

import (
	"log"

	"github.com/spf13/viper"
)

var cfg Config

type Config struct {
	reader *viper.Viper
}

func new() Config {
	reader := viper.New()
	reader.AddConfigPath("./")
	reader.SetConfigFile("config.yml")

	if err := reader.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	return Config{reader}
}

// Return global instance
func GetInstance() *Config {
	if !cfg.isInstantiated() {
		cfg = new()
	}

	return &cfg
}

func (c *Config) readEnv(key string) string {
	return c.reader.GetString(key)
}

func (c *Config) FetchDatabaseConfig() *DatabaseConfig {

	port := c.reader.GetInt("Development.Database.Port")

	return &DatabaseConfig{
		User:         c.readEnv("Development.Database.User"),
		Password:     c.readEnv("Development.Database.Password"),
		Host:         c.readEnv("Development.Database.Host"),
		DatabaseName: c.readEnv("Development.Database.Name"),
		Port:         port,
	}

}

func (c *Config) FetchServerPort() uint64 {
	return c.reader.GetUint64("Server.Port")
}

func (c *Config) GetJWTSecret() string {
	return c.reader.GetString("JWT.Secret")
}

func (c *Config) GetBucketName() string {
	return c.reader.GetString("Oracle.BucketName")
}

func (c *Config) isInstantiated() bool {
	return c.reader != nil
}
