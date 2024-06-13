package util

import (
	"backend/server/env"
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type Encryption struct{}

func (ref Encryption) EncodeWithSecret(text string) (string, error) {
	secretKey := env.GetSecretKey()
	return ref.Encode(secretKey, text)
}
func (ref Encryption) DecodeWithSecret(encodedText string) (string, error) {
	secretKey := env.GetSecretKey()
	return ref.Decode(secretKey, encodedText)
}

func (ref Encryption) Encode(secretKey string, text string) (string, error) {
	key := hashKey(secretKey)

	// Layer 1: DES Encryption with original hashed key
	cipherText, err := encryptMethod(key, text)
	if err != nil {
		return "", err
	}

	// Layer 2: DES Encryption with reversed hashed key
	reversedKey := hashKey(reverseStrings(string(key)))
	cipherText, err = encryptMethod(reversedKey, cipherText)
	if err != nil {
		return "", err
	}

	// Layer 3: DES Encryption with first half of the original hashed key rehashed
	firstHalfKey := hashKey(string(key)[:len(key)/2])
	cipherText, err = encryptMethod(firstHalfKey, cipherText)
	if err != nil {
		return "", err
	}

	// Layer 4: DES Encryption with second half of the original hashed key rehashed
	secondHalfKey := hashKey(string(key)[len(key)/2:])
	cipherText, err = encryptMethod(secondHalfKey, cipherText)
	if err != nil {
		return "", err
	}

	// Layer 5: Base64 Encoding
	return base64.StdEncoding.EncodeToString([]byte(cipherText)), nil
}

func (ref Encryption) Decode(secretKey string, encodedText string) (string, error) {
	key := hashKey(secretKey)

	// Layer 5: Base64 Decoding
	cipherTextBytes, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}

	// Layer 4: DES Decryption with second half of the original hashed key rehashed
	plaintext, err := decryptMethod(hashKey(string(key)[len(key)/2:]), string(cipherTextBytes))
	if err != nil {
		return "", err
	}

	// Layer 3: DES Decryption with first half of the original hashed key rehashed
	plaintext, err = decryptMethod(hashKey(string(key)[:len(key)/2]), plaintext)
	if err != nil {
		return "", err
	}

	// Layer 2: DES Decryption with reversed hashed key
	reversedKey := hashKey(reverseStrings(string(key)))
	plaintext, err = decryptMethod(reversedKey, plaintext)
	if err != nil {
		return "", err
	}

	// Layer 1: DES Decryption with original hashed key
	plaintext, err = decryptMethod(key, plaintext)
	if err != nil {
		return "", err
	}

	return plaintext, nil
}

func hashKey(key string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	fullHash := hasher.Sum(nil)
	// Truncate the hash to 24 bytes for Triple DES
	return fullHash[:24]
}

func reverseStrings(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	padding := int(src[length-1])
	return src[:length-padding], nil
}

func encryptMethod(key []byte, plaintext string) (string, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}

	plaintextBytes := pad([]byte(plaintext), block.BlockSize())
	cipherText := make([]byte, des.BlockSize+len(plaintextBytes))
	iv := cipherText[:des.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[des.BlockSize:], plaintextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decryptMethod(key []byte, cipherText string) (string, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < des.BlockSize || len(cipherTextBytes)%des.BlockSize != 0 {
		return "", fmt.Errorf("cipher text length must be a multiple of block size")
	}

	iv := cipherTextBytes[:des.BlockSize]
	cipherTextBytes = cipherTextBytes[des.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)

	plaintextBytes, err := unpad(cipherTextBytes)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}
