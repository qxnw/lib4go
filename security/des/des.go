package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Encrypt DES加密
// input 要加密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func Encrypt(input string, skey string) (r string, err error) {
	origData := []byte(input)
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("des NewCipher err:%v", err)
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

// Decrypt DES解密
// input 要解密的字符串	skey 加密使用的秘钥[字符串长度必须是8的倍数]
func Decrypt(input string, skey string) (r string, err error) {
	/*add by champly 2016年11月16日17:35:03*/
	if len(input) < 1 {
		return "", errors.New("解密的对象长度必须大于0")
	}
	/*end*/

	crypted, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("hex DecodeString err:%v", err)
	}
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("des NewCipher err:%v", err)
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	r = string(origData)
	return
}

// ZeroPadding Zero填充模式
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// ZeroUnPadding 去除Zero的补码
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

// PKCS5Padding PKCS5填充模式
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding 去除PKCS5的补码
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// PKCS7Padding PKCS7填充模式
func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

// PKCS7UnPadding 去除PKCS7的补码
func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
