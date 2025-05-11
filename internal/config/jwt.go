package config

import (
	"os"
	"time"
)

// JWTConfig содержит настройки для JWT-токенов
type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func LoadJWT() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "zrvbJBByrUEDl4994VA4vooeVCimuVlIJjWxJ2ZZrnS7aCIYEOar6ExAH3dDhEFT"
	}

	return JWTConfig{
		Secret:    secret,
		ExpiresIn: 24 * time.Hour,
	}
}
