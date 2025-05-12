package config

import "github.com/sirupsen/logrus"

type CryptoConfig struct {
	PGPKey  string
	HMACKey string
}

func LoadCrypto() CryptoConfig {
	logrus.Info("Загружаем конфиг для шифрования")
	return CryptoConfig{
		PGPKey:  getEnv("PGP_KEY", "PGP_KEY"),
		HMACKey: getEnv("HMAC_KEY", "HMAC_KEY"),
	}
}
