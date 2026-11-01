package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func genKey(password string) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	return h.Sum(nil)
}

func genIV(password string) []byte {
	iv := make([]byte, 16)
	for i := 0; i < 16 && i < len(password); i++ {
		iv[i] = password[i]
	}
	return iv
}

func Encrypt(data string, password string) ([]byte, error) {
	hash := genKey(password)
	// iv := genIV(password)
	plaintText := []byte(data)

	block, err := aes.NewCipher(hash)
	if err != nil {
		return []byte("Error"), fmt.Errorf("error cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte("Error"), fmt.Errorf("error gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte("Error"), fmt.Errorf("error read: %w", err)
	}

	return gcm.Seal(nonce, nonce, plaintText, nil), nil
}

func Decrypt(data []byte, password string) (string, error) {
	hash := genKey(password)

	block, err := aes.NewCipher(hash)
	if err != nil {
		return "Error", fmt.Errorf("error cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "Error", fmt.Errorf("error gcm: %w", err)
	}

	nonseSize := gcm.NonceSize()
	if len(data) < nonseSize {
		return "Error", fmt.Errorf("encrypted data is too short")
	}

	nonce, cipherText := data[:nonseSize], data[nonseSize:]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "Error", fmt.Errorf("error decrypt: %w", err)
	}

	return string(plainText), nil
}
