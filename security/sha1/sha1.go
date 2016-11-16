package sha1

import (
	"crypto/sha1"
	"fmt"
)

func Encrypt(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
