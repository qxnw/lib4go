package des

import (
	"strings"
	"testing"
)

func TestEncrypt(t *testing.T) {
	input := "hello"
	key := "12345678"
	except := "BA16C6A0257125AF"
	actual, err := Encrypt(input, key)
	if err != nil {
		t.Errorf("Encrypt fail %v", err)
	}
	if !strings.EqualFold(except, actual) {
		t.Errorf("Encrypt fail %s to %s", except, actual)
	}

	input = "hello"
	key = "123456789"
	_, err = Encrypt(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = ""
	key = "12345678"
	except = "FEB959B7D4642FCB"
	actual, err = Encrypt(input, key)
	if err != nil {
		t.Error("test fail")
	}
	if !strings.EqualFold(except, actual) {
		t.Errorf("Encrypt fail %s to %s", except, actual)
	}
}

func TestDecrypt(t *testing.T) {
	input := "BA16C6A0257125AF"
	key := "12345678"
	except := "hello"
	actual, err := Decrypt(input, key)
	if err != nil {
		t.Errorf("Decrypt fail %v", err)
	}
	if !strings.EqualFold(except, actual) {
		t.Errorf("Decrypt fail %s to %s", except, actual)
	}

	input = "BA16C6A0257125AF"
	key = "123456789"
	_, err = Decrypt(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = ""
	key = "12345678"
	_, err = Decrypt(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = "0"
	key = "12345678"
	_, err = Decrypt(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = "!#@@@#!"
	key = "12345678"
	_, err = Decrypt(input, key)
	if err == nil {
		t.Error("test fail")
	}
}

func TestEncrypt3DES(t *testing.T) {
	input := "hello"
	key := "123456781234567812345678"
	except := "BA16C6A0257125AF"
	actual, err := Encrypt3DES(input, key)
	if err != nil {
		t.Errorf("Encrypt fail %v", err)
	}
	if !strings.EqualFold(except, actual) {
		t.Errorf("Encrypt fail %s to %s", except, actual)
	}

	input = "hello"
	key = "123456789"
	_, err = Encrypt3DES(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = ""
	key = "123456781234567812345678"
	except = "FEB959B7D4642FCB"
	actual, err = Encrypt3DES(input, key)
	if err != nil {
		t.Error("test fail")
	}
	if !strings.EqualFold(except, actual) {
		t.Errorf("Encrypt fail %s to %s", except, actual)
	}
}

func TestDecrypt3DES(t *testing.T) {
	input := "BA16C6A0257125AF"
	key := "123456781234567812345678"
	except := "hello"
	actual, err := Decrypt3DES(input, key)
	if err != nil {
		t.Errorf("Decrypt3DES fail %v", err)
	}
	if !strings.EqualFold(except, actual) {
		t.Errorf("Decrypt3DES fail %s to %s", except, actual)
	}

	input = "BA16C6A0257125AF"
	key = "12345678"
	_, err = Decrypt3DES(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = ""
	key = "123456781234567812345678"
	_, err = Decrypt3DES(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = "0"
	key = "123456781234567812345678"
	_, err = Decrypt3DES(input, key)
	if err == nil {
		t.Error("test fail")
	}

	input = "!#@@@#!"
	key = "123456781234567812345678"
	_, err = Decrypt3DES(input, key)
	if err == nil {
		t.Error("test fail")
	}
}
