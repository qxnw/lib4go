package rsa

/*
	http://blog.studygolang.com/2013/01/go%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86%E4%B9%8Brsa/
	2016年11月16日21:34:36
*/

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/qxnw/lib4go/encoding"
)

var (
	privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDTeyBioVCJcwhwu68CKJQQwrHgzRu+14fM5/Br/xBrpp1+FPtB
MKSs+jlJ5gMjVOp1gYCRvQ2pSdaxj83t7i0KRRzkUCup6rIphC6265XGWagCODT2
LaUkTnIBA3wrbDF7ziTrHJt0e6WO4NJie4iGzqfKbQMbXnj9dvSlV3KWIQIDAQAB
AoGBAMnRAJDfPQpOesmKcnLu4o40HqhXVJkE+hWzah7F5Je3AyklQLlvgFeK200I
cgovqSfGFDoAXp8lVftRLsZWuycKykhqwZHD7WVN0b3Jtnxr9Q+ZB3vpnu+/sRGF
9EJ0t0Q0ESExMgRkaDeDiENsrn8KZ0EpXwUo2IlpxwCGKp65AkEA+IOyp7ytsg4f
zXghTnR+tRzSMu2ou4pPx0bfXE1qkHLdw/6D9xYlPmV8eHFEQpQO/kAny9D4mHA1
UOj1QXII7wJBANnZ3MSUKoriwgLqxEuh3VjxiXlgAWnGkbL0aG3pKIMHVZE630g3
5lHBRYResvm5BukPtLPxMsdag8gpQ7Ehse8CQQDjPdHkjaQqt72e7aVPDzk5tWQE
C8uJybyPlR/zUBsMgOyGJrpW+xoNR1Gc9L2dP7PCC7oYJjrbcWdfV9XEBVljAkA/
Gernid9UwV/fBm97VNRPmg7u+E8Qe3LiegbxpzKT2YEAgyP/wClXjvr6349J5D1L
LsBxyrChq+c2CDXSTedDAkAdkrx++tGeitQOYjP8EdC3roIq1x5+eAm78HRQ+SY2
SGjbwWCzZ/26kW0bm6H5KcOS1xPjp6gn+Bd00PQRaGJ5
-----END RSA PRIVATE KEY-----`)
	publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDTeyBioVCJcwhwu68CKJQQwrHg
zRu+14fM5/Br/xBrpp1+FPtBMKSs+jlJ5gMjVOp1gYCRvQ2pSdaxj83t7i0KRRzk
UCup6rIphC6265XGWagCODT2LaUkTnIBA3wrbDF7ziTrHJt0e6WO4NJie4iGzqfK
bQMbXnj9dvSlV3KWIQIDAQAB
-----END PUBLIC KEY-----`)
)

func TestEncrypt(t *testing.T) {
	input := "hello"
	except := "TkNIRnIZESmOdid6j5ObIeTKBUmHMhqzmYp6A2k/8TtKOOmheBv2Ji2ufDxHiC+7KdwKaaWdMnAKXvGuZ1QrP6Q9+i4c8MWSFBfCDuiQXqH6lLXer6k4pq2LUH9TIXg1HQB38Kn3eWElQW7AN/IKdLpM2VMSIy4Rd3SnEOh62ZA="
	data, err := Encrypt(input, string(publicKey))
	fmt.Println(bytes.EqualFold(data, []byte(except)))
	if err != nil {
		t.Errorf("encrypt fail:%v", err)
	}
	if !bytes.EqualFold(data, []byte(except)) {
		t.Errorf("encrypt fail %s to %s", except, string(data))
	}
}

func TestDecode(t *testing.T) {
	input := "TkNIRnIZESmOdid6j5ObIeTKBUmHMhqzmYp6A2k/8TtKOOmheBv2Ji2ufDxHiC+7KdwKaaWdMnAKXvGuZ1QrP6Q9+i4c8MWSFBfCDuiQXqH6lLXer6k4pq2LUH9TIXg1HQB38Kn3eWElQW7AN/IKdLpM2VMSIy4Rd3SnEOh62ZA="
	except := "hello"
	data, err := Decrypt(input, string(privateKey))
	if err != nil {
		t.Errorf("Decode fail:%v", err)
	}
	if !bytes.EqualFold(data, []byte(except)) {
		t.Errorf("Decode fail %s to %s", except, string(data))
	}
}

func TestSign(t *testing.T) {
	input := "hello"
	except := ""
	data, err := Sign(input, string(publicKey), "md5")
	if err != nil {
		t.Errorf("Sign fail %v", err)
	}
	if !bytes.EqualFold(data, []byte(except)) {
		t.Errorf("Sign fail %s to %s", except, string(data))
	}

	input = "hello"
	except = ""
	data, err = Sign(input, string(publicKey), "sha1")
	if err != nil {
		t.Errorf("Sign fail %v", err)
	}
	if !bytes.EqualFold(data, []byte(except)) {
		actual, _ := encoding.Convert(data, "gb2312")
		t.Errorf("Sign fail %s to %s", except, actual)
	}

	input = "hello"
	except = ""
	data, err = Sign(input, string(publicKey), "base64")
	if err == nil {
		t.Error("test fail")
	}

	input = "hello"
	except = ""
	data, err = Sign(input, string(publicKey), "")
	if err == nil {
		t.Error("test fail")
	}
}

func TestVerify(t *testing.T) {
	input := "hello"
	sign := ""
	flag, err := Verify(input, sign, string(privateKey), "md5")
	if err != nil {
		t.Errorf("Verify fail %v", err)
	}
	if flag {
		t.Error("Verify fail")
	}

	input = "hello"
	sign = ""
	flag, err = Verify(input, sign, string(privateKey), "sha1")
	if err != nil {
		t.Errorf("Verify fail %v", err)
	}
	if flag {
		t.Error("Verify fail")
	}

	input = "hello"
	sign = ""
	flag, err = Verify(input, sign, string(privateKey), "base64")
	if err == nil {
		t.Error("Verify fail")
	}

	input = "hello"
	sign = ""
	flag, err = Verify(input, sign, string(privateKey), "")
	if err == nil {
		t.Error("Verify fail")
	}
}
