package config

import (
	"github.com/sirupsen/logrus"
)

// CryptoConfig содержит ключи для шифрования и подписи
type CryptoConfig struct {
	PGPKey  string // Ключ для PGP-шифрования
	HMACKey string // Ключ для HMAC-подписей
}

// LoadCrypto загружает конфигурацию криптографических ключей
func LoadCrypto() CryptoConfig {
	cfg := CryptoConfig{
		PGPKey:  getEnv("BANK_PGP_KEY", "PRaRpVX3wyGwNSH"),
		HMACKey: getEnv("BANK_HMAC_KEY", "5MEMgsKNAMnaa0W"),
	}

	// Логирование (без вывода самих ключей)
	logrus.Info("Конфигурация криптографических ключей загружена")

	return cfg
}
