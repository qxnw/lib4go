package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strings"

	"github.com/qxnw/lib4go/encoding/base64"
	"github.com/qxnw/lib4go/security/des"
)

func getKey(key string) []byte {
	arrKey := []byte(key)
	keyLen := len(key)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

//Encrypt 加密字符串
func Encrypt(msg string, key string) (string, error) {
	keyBytes := getKey(key)
	var iv = keyBytes[:aes.BlockSize]
	encrypted := make([]byte, len(msg))
	aesBlockEncrypter, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(msg))
	return base64.Encode(string(encrypted)), nil
}

//Decrypt 解密字符串
func Decrypt(src string, key string) (msg string, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	content, err := base64.Decode(src)
	if err != nil {
		return
	}
	keyBytes := getKey(key)
	var iv = keyBytes[:aes.BlockSize]
	decrypted := make([]byte, len(content))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, []byte(content))
	return strings.TrimSpace(string(decrypted)), nil
}

// EncryptCBCPKCS7 CBC模式,PKCS7填充
func EncryptCBCPKCS7(contentStr string, keyStr string, iv []byte) (string, error) {
	content := []byte(contentStr)
	key := []byte(keyStr)
	if len(content)&aes.BlockSize != 0 {
		return "", fmt.Errorf("要加密的明文不是块大小的倍数")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(iv) != block.BlockSize() {
		return "", fmt.Errorf("IV length must equal block size")
	}

	content = des.PKCS7Padding(content)
	// iv := key[:block.BlockSize()]
	blockModel := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(content))
	blockModel.CryptBlocks(cipherText, content)
	return base64.EncodeBytes(cipherText), nil
}

// DecryptCBCPKCS7 CBC模式,PKCS7填充
func DecryptCBCPKCS7(contentStr string, keyStr string, iv []byte) (string, error) {
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

	if len(iv) != block.BlockSize() {
		return "", fmt.Errorf("IV length must equal block size")
	}

	// iv := key[:block.BlockSize()]
	blockModel := cipher.NewCBCDecrypter(block, iv)

	plantText := make([]byte, len(content))
	blockModel.CryptBlocks(plantText, content)
	plantText = des.UnPKCS7Padding(plantText)

	return string(plantText), nil
}
