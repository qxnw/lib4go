package hex

import "encoding/hex"

func Encode(src []byte) string {
	return hex.EncodeToString(src)
}

func Decode(src string) (r string, err error) {
	data, err := hex.DecodeString(src)
	if err != nil {
		return
	}
	r = string(data)
	return
}