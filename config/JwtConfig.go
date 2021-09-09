package config

type UserJwt struct {
	ID   uint64
	Ic   string
	Role []string
}

func (c *Config) GetJWTSecret() string {
	return CFG.reader.GetString("JWT.Secret")
}
