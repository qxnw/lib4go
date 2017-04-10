package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/qxnw/lib4go/encoding/base64"
	"github.com/qxnw/lib4go/security/des"
)

// EncryptCBCPKCS7 CBC模式,PKCS7填充
func EncryptCBCPKCS7(contentStr string, keyStr string) (string, error) {
	content := []byte(contentStr)
	key := []byte(keyStr)
	if len(content)&aes.BlockSize != 0 {
		return "", fmt.Errorf("要加密的明文不是块大小的倍数")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	content = des.PKCS7Padding(content)
	iv := make([]byte, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(content))
	blockModel.CryptBlocks(cipherText, content)
	return base64.EncodeBytes(cipherText), nil
}

// DecryptCBCPKCS7 CBC模式,PKCS7填充
func DecryptCBCPKCS7(contentStr string, keyStr string) (string, error) {
	content, err := base64.DecodeBytes(contentStr)
	if err != nil {
		return "", err
	}

	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(content) < aes.BlockSize {
		return "", fmt.Errorf("要解密的字符串太短")
	}

	iv := make([]byte, block.BlockSize())
	blockModel := cipher.NewCBCDecrypter(block, iv)

	plantText := make([]byte, len(content))
	blockModel.CryptBlocks(plantText, content)
	plantText = des.UnPKCS7Padding(plantText)

	return string(plantText), nil
}
