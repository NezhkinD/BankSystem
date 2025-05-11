package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var encryptionKey = []byte("32_byte_key_here_32_byte_key_here")

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := []byte("unique_nonce123") // В реальности генерируйте уникальный nonce
	encrypted := gcm.Seal(nil, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
