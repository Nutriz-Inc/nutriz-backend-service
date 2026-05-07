package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"nutriz-backend-service/config"

	"github.com/google/uuid"
)

type Secret struct{}

func (Secret) Generate() string {
	sum := md5.Sum([]byte(uuid.NewString()))
	return hex.EncodeToString(sum[:])
}

func (s Secret) Encrypt(cfg *config.Env, secret string) (string, error) {
	block, iv, err := s.getCipher(cfg)
	if err != nil {
		return "", err
	}

	plain := s.pkcs7Padding([]byte(secret), aes.BlockSize)

	encrypted := make([]byte, len(plain))

	cipher.NewCBCEncrypter(block, iv).
		CryptBlocks(encrypted, plain)

	return hex.EncodeToString(encrypted), nil
}

func (s Secret) Decrypt(cfg *config.Env, encryptedHex string) (string, error) {
	block, iv, err := s.getCipher(cfg)
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", fmt.Errorf("decode encrypted data: %w", err)
	}

	if len(encrypted) == 0 {
		return "", fmt.Errorf("encrypted data is empty")
	}

	if len(encrypted)%aes.BlockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext block size")
	}

	decrypted := make([]byte, len(encrypted))

	cipher.NewCBCDecrypter(block, iv).
		CryptBlocks(decrypted, encrypted)

	decrypted, err = s.pkcs7Unpadding(decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (s Secret) IsEqual(cfg *config.Env, incoming string, encrypted string) bool {
	decrypted, err := s.Decrypt(cfg, encrypted)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(
		[]byte(incoming),
		[]byte(decrypted),
	) == 1
}

func (s Secret) getCipher(cfg *config.Env) (cipher.Block, []byte, error) {
	key, err := hex.DecodeString(cfg.Secret.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("decode key: %w", err)
	}

	iv, err := hex.DecodeString(cfg.Secret.IV)
	if err != nil {
		return nil, nil, fmt.Errorf("decode iv: %w", err)
	}

	switch len(key) {
	case 16, 24, 32:
	default:
		return nil, nil, fmt.Errorf("invalid AES key size: %d", len(key))
	}

	if len(iv) != aes.BlockSize {
		return nil, nil, fmt.Errorf("invalid IV size: %d", len(iv))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("create cipher: %w", err)
	}

	return block, iv, nil
}

func (Secret) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize

	return append(
		data,
		bytes.Repeat([]byte{byte(padding)}, padding)...,
	)
}

func (Secret) pkcs7Unpadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty decrypted data")
	}

	padding := int(data[len(data)-1])

	if padding == 0 || padding > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}

	for _, b := range data[len(data)-padding:] {
		if int(b) != padding {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
