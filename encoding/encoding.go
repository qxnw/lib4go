package encoding

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetReader(content string, charset string) io.Reader {
	if strings.EqualFold(charset, "utf-8") {
		return strings.NewReader(content)
	}
	return transform.NewReader(bytes.NewReader([]byte(content)), simplifiedchinese.GBK.NewDecoder())
}

func Convert(buf []byte, charset string) string {
	if strings.EqualFold(charset, "utf-8") {
		return string(buf)
	}
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(buf), simplifiedchinese.GBK.NewDecoder()))
	if err == nil {
		return string(data)
	}
	return string(buf)
}

func UnicodeEncode(str string) string {
	rs := []rune(str)
	json := ""
	html := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
			html += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}
	return json
}

func UnicodeDecode(unicodeStr string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(unicodeStr)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(unicodeStr[i:])
			break
		}
		if unicodeStr[i] == '\\' && unicodeStr[i+1] == 'u' {
			hex := unicodeStr[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(unicodeStr[i:x])
			}
			i = x
		} else {
			buf.WriteByte(unicodeStr[i])
			i++
		}
	}
	return buf.String()
}
