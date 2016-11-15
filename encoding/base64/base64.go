package base64

import "encoding/base64"

func EncodeBytes(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Encode(src string) string {
	return EncodeBytes([]byte(src))
}

func DecodeBytes(src string) (s []byte, err error) {
	s, err = base64.StdEncoding.DecodeString(src)
	return
}
func Decode(src string) (s string, err error) {
	buf, err := DecodeBytes(src)
	if err != nil {
		return
	}
	s = string(buf)
	return
}
