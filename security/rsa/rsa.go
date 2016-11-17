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

/*
	RSA加密
	params
		origData 要加密的参数
		publicKey 加密时候用到的公钥
	return
		string 加密之后的字符串
		error 加密时产生的错误
*/
func Encrypt(origData string, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	// return rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))

	/*change by champly 2016年11月17日09:41:21*/
	data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
	/*end*/
}

/*
	RSA解密
	params
		ciphertext 要解密的参数
		privateKey 解密时候用到的公钥
	return
		string 解密之后的字符串
		error 解密时产生的错误
*/
func Decrypt(ciphertext string, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// return rsa.DecryptPKCS1v15(rand.Reader, priv, []byte(ciphertext))

	/*change by champly 2016年11月17日09:36:41*/
	input, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", nil
	}
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, input)
	if err != nil {
		return "", err
	}

	return string(data), nil
	/*end*/
}

/*
	使用RSA生成签名
	params
		message 要签名的字符串
		privateKey 加密时使用的秘钥
		mode 加密的模式【目前只支持MD5，SHA1，不区分大小写】
	return
		string 签名之后的字符串
		error 签名时产生成错误信息
*/
func Sign(message string, privateKey string, mode string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	switch strings.ToLower(mode) {
	case "sha1":
		t := sha1.New()
		io.WriteString(t, message)
		digest := t.Sum(nil)
		// return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, digest)
		/*change by champly 2016年11月17日10:31:10*/
		data, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, digest)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
		/*end*/
	case "md5":
		t := md5.New()
		io.WriteString(t, message)
		digest := t.Sum(nil)
		// return rsa.SignPKCS1v15(rand.Reader, priv, crypto.MD5, digest)
		/*change by champly 2016年11月17日10:31:10*/
		data, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.MD5, digest)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
		/*end*/
	default:
		return "", errors.New("签名模式不支持")
	}

}

/*
	Verify校验签名
	params
		src	要验证的签名字符串
		sign 生成的签名字符串
		publicKey 验证签名的公钥
		mode 加密的模式【目前只支持MD5，SHA1，不区分大小写】
	return
		pass 是否通过校验
		err 校验的时候产生的错误
*/
func Verify(src string, sign string, publicKey string, mode string) (pass bool, err error) {
	//步骤1，加载RSA的公钥
	block, _ := pem.Decode([]byte(publicKey))
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
