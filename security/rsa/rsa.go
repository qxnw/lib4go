package rsa

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"strings"
)

//Encrypt RSA加密
func Encrypt(origData string, publicKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
}

//Decrypt RSA解密
func Decrypt(ciphertext string, privateKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(ciphertext))
}

//Sign 生成签名
func Sign(message string, privateKey string, mode string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(mode) {
	case "sha1":
		t := sha1.New()
		io.WriteString(t, message)
		digest := t.Sum(nil)
		return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, digest)
	case "md5":
		t := md5.New()
		io.WriteString(t, message)
		digest := t.Sum(nil)
		return rsa.SignPKCS1v15(rand.Reader, priv, crypto.MD5, digest)
	default:
		return nil, errors.New("签名模式不支持")
	}

}

//Verify 验签
func Verify(src string, sign string, pubkey string, mode string) (pass bool, err error) {
	//步骤1，加载RSA的公钥
	block, _ := pem.Decode([]byte(pubkey))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	rsaPub, _ := pub.(*rsa.PublicKey)
	data, _ := base64.StdEncoding.DecodeString(sign)
	switch strings.ToLower(mode) {
	case "sha1":
		t := sha1.New()
		io.WriteString(t, src)
		digest := t.Sum(nil)
		err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA1, digest, data)
	case "md5":
		t := md5.New()
		io.WriteString(t, src)
		digest := t.Sum(nil)
		err = rsa.VerifyPKCS1v15(rsaPub, crypto.MD5, digest, data)
	default:
		err = errors.New("验签模式不支持")
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
