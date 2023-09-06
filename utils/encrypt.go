package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
)

// Encrypt 使用AES加密
func Encrypt(plaintext, key string) (string, error) {
	plaintextByte := []byte(plaintext)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintextByte))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextByte)

	logx.Info(fmt.Printf("Original: %s\n", plaintext))
	logx.Info(fmt.Printf("Encrypted: %x\n", ciphertext))

	return hex.EncodeToString(ciphertext), nil

}

// Decrypt 使用AES解密
func Decrypt(ciphertext, key string) (string, error) {
	keyByte := []byte(key)

	cipherTextByte, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	if len(cipherTextByte) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := cipherTextByte[:aes.BlockSize]
	cipherTextByte = cipherTextByte[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherTextByte, cipherTextByte)

	logx.Info(fmt.Printf("Decrypted: %s\n", cipherTextByte))

	return string(cipherTextByte), nil
}
