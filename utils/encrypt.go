package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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

	fmt.Printf("Original: %s", plaintext)
	fmt.Println()
	fmt.Printf("Encrypted: %x", ciphertext)
	fmt.Println()
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
		return "", errors.New("ciphertext too short")
	}

	iv := cipherTextByte[:aes.BlockSize]
	cipherTextByte = cipherTextByte[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherTextByte, cipherTextByte)

	fmt.Println()
	fmt.Printf("Decrypted: %s", cipherTextByte)
	fmt.Println()

	return string(cipherTextByte), nil
}
