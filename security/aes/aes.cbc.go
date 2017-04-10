package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func EncryptCBC(msg, keyStr string) ([]byte, error) {
	key := []byte(keyStr)
	plantText := []byte(msg)
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = PKCS7Padding(plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, plantText)
	ciphertext := make([]byte, len(plantText))
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func DecryptCBC(msg, keyStr string) ([]byte, error) {
	key := []byte(keyStr)
	ciphertext := []byte(msg)
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, ciphertext)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = PKCS7UnPadding(ciphertext, block.BlockSize())
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
