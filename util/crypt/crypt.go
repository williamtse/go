package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"os"
)

// 加密函数
func CryptoEncrypt(plainData []byte) string {
	encodedKey := os.Getenv("CRYPTO_KEY") // 32字节密钥
	encodediv := os.Getenv("CRYPTO_IV")   // 16字节初始向量

	key := []byte(encodedKey)
	iv := []byte(encodediv)

	// // 解码 base64 密钥
	// key, err := base64.StdEncoding.DecodeString(encodedKey)
	// if err != nil {
	// 	log.Fatal("Error decoding key:", err)
	// }

	// // 解码 base64 密钥
	// iv, err := base64.StdEncoding.DecodeString(encodediv)
	// if err != nil {
	// 	log.Fatal("Error decoding iv:", err)
	// }

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	paddedData := pkcs7Pad(plainData, aes.BlockSize)

	cipherText := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddedData)

	return base64.StdEncoding.EncodeToString(cipherText)
}

// PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 解密函数
func CryptoDecrypt(encryptedData string) (string, error) {
	encodedKey := os.Getenv("CRYPTO_KEY") // 32字节密钥
	encodediv := os.Getenv("CRYPTO_IV")   // 16字节初始向量
	key := []byte(encodedKey)
	iv := []byte(encodediv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	return string(pkcs7Unpad(cipherText)), nil
}

// PKCS7去除填充
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
