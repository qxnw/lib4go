package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"strings"
)

/*
   DES加密
   params:
       input 要加密的字符串
       skey  加密使用的秘钥【字符串的长度必须是8的倍数】
   return
        r    加密后的结果
        err  加密的时候出现的异常
*/
func Encrypt(input string, skey string) (r string, err error) {
	origData := []byte(input)
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

/*
   DES解密
   params:
       input 要解密的字符串
       skey  解密使用的秘钥【字符串的长度必须是8的倍数】
   return
        r    解密后的结果
        err  解密的时候出现的异常
*/
func Decrypt(input string, skey string) (r string, err error) {
	/*add by champly 2016年11月16日17:35:03*/
	if len(input) < 1 {
		return r, errors.New("解密的对象长度必须大于0")
	}
	/*end*/

	crypted, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	r = string(origData)
	return
}

/*
   3DES加密
   params:
       input 要加密的字符串
       skey  加密使用的秘钥【字符串的长度必须是24的倍数】
   return
        r    加密后的结果
        err  加密的时候出现的异常
*/
func Encrypt3DES(input string, skey string) (r string, err error) {
	origData := []byte(input)
	key := []byte(skey)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

/*
   3DES解密
   params:
       input 要解密的字符串
       skey  解密使用的秘钥【字符串的长度必须是24的倍数】
   return
        r    解密后的结果
        err  解密的时候出现的异常
*/
func Decrypt3DES(input, skey string) (r string, err error) {
	/*add by champly 2016年11月16日17:35:03*/
	if len(input) < 1 {
		return r, errors.New("解密的对象长度必须大于0")
	}
	/*end*/

	crypted, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	key := []byte(skey)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	r = string(origData)
	return
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

/**
 *  去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
