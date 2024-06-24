package util

import (
	"backend/server/env"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

// -------------------------------------------------------------------
// -------------------------------------------------------------------

func hashKey(key string, length int) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	hashed := hasher.Sum(nil)
	return hashed[:length]
}

func reverseString(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// -------------------------------------------------------------------
// -------------------------------------------------------------------

func encryptTripleDES(plaintext string, secret_key string) (string, error) {
	key := hashKey(secret_key, 24)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}

	// Padding plaintext agar panjangnya menjadi kelipatan 8 byte (ukuran blok DES)
	padding := des.BlockSize - (len(plaintext) % des.BlockSize)
	padtext := []byte(plaintext + strings.Repeat(string(byte(padding)), padding))

	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, key[:des.BlockSize])
	mode.CryptBlocks(ciphertext, padtext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptTripleDES(ciphertext string, secret_key string) (string, error) {
	key := hashKey(secret_key, 24)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < des.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	mode := cipher.NewCBCDecrypter(block, key[:des.BlockSize])
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	// Unpad plaintext by removing padding bytes
	padding := int(ciphertextBytes[len(ciphertextBytes)-1])
	return string(ciphertextBytes[:len(ciphertextBytes)-padding]), nil
}

// -------------------------------------------------------------------

func encryptAES(plaintext string, secret_key string) (string, error) {
	key := hashKey(secret_key, 32)
	iv := hashKey(secret_key, 16)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Padding plaintext agar panjangnya menjadi kelipatan 16 byte (ukuran blok AES)
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	padtext := []byte(plaintext + strings.Repeat(string(byte(padding)), padding))

	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padtext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(ciphertext string, secret_key string) (string, error) {
	key := hashKey(secret_key, 32)
	iv := hashKey(secret_key, 16)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		panic("ciphertext too short")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	// Unpad plaintext by removing padding bytes
	padding := int(ciphertextBytes[len(ciphertextBytes)-1])
	return string(ciphertextBytes[:len(ciphertextBytes)-padding]), nil
}

// -------------------------------------------------------------------

func encryptMethod(plainText, key string) (string, error) {
	// fmt.Println("plainText:", plainText, "key:", key)

	return encryptTripleDES(plainText, key)
	// return encryptAES(plainText, key)
}
func decryptMethod(cipherText, key string) (string, error) {
	return decryptTripleDES(cipherText, key)
	// return decryptAES(cipherText, key)
}

// -------------------------------------------------------------------
// -------------------------------------------------------------------

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
	// Layer 1: AES Encryption with original hashed key
	cipherText, err := encryptMethod(text, secretKey)
	if err != nil {
		return "", err
	}
	log.Println("GO Encode Layer 1:", cipherText)

	// Layer 2: AES Encryption with reversed hashed key
	reversedKey := reverseString(secretKey)
	cipherText, err = encryptMethod(cipherText, reversedKey)
	if err != nil {
		return "", err
	}
	log.Println("GO Encode Layer 2:", cipherText)

	// Layer 3: AES Encryption with first half of the original hashed key rehashed
	firstHalfKey := string(secretKey)[:len(secretKey)/2]
	cipherText, err = encryptMethod(cipherText, firstHalfKey)
	if err != nil {
		return "", err
	}
	log.Println("GO Encode Layer 3:", cipherText)

	// Layer 4: AES Encryption with second half of the original hashed key rehashed
	secondHalfKey := string(secretKey)[len(secretKey)/2:]
	cipherText, err = encryptMethod(cipherText, secondHalfKey)
	if err != nil {
		return "", err
	}
	log.Println("GO Encode Layer 4:", cipherText)

	// Layer 5: Base64 Encoding
	cipherText = base64.StdEncoding.EncodeToString([]byte(cipherText))
	log.Println("GO Encode Layer 5:", cipherText)
	return cipherText, nil
}

func (ref Encryption) Decode(secretKey string, encodedText string) (string, error) {
	// Layer 5: Base64 Decoding
	cipherTextBytes, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}
	log.Println("GO Decode Layer 5:", string(cipherTextBytes))

	// Layer 4: AES Decryption with second half of the original hashed key rehashed
	plaintext, err := decryptMethod(string(cipherTextBytes), string(secretKey)[len(secretKey)/2:])
	if err != nil {
		return "", err
	}
	log.Println("GO Decode Layer 4:", plaintext)

	// Layer 3: AES Decryption with first half of the original hashed key rehashed
	plaintext, err = decryptMethod(plaintext, string(secretKey)[:len(secretKey)/2])
	if err != nil {
		return "", err
	}
	log.Println("GO Decode Layer 3:", plaintext)

	// Layer 2: AES Decryption with reversed hashed key
	reversedKey := reverseString(secretKey)
	plaintext, err = decryptMethod(plaintext, reversedKey)
	if err != nil {
		return "", err
	}
	log.Println("GO Decode Layer 2:", plaintext)

	// Layer 1: AES Decryption with original hashed key
	plaintext, err = decryptMethod(plaintext, secretKey)
	if err != nil {
		return "", err
	}
	log.Println("GO Decode Layer 1:", plaintext)

	return plaintext, nil
}

// ------------------------------------------------------
// ------------------------------------------------------
