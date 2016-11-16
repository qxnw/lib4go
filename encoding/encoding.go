package encoding

import (
	"bytes"
	"io"
	"io/ioutil"
<<<<<<< HEAD
	"strings"

	"fmt"

=======
	"strconv"
	"strings"

>>>>>>> e207abb0f14858dd24a18e962477dca1545b8cec
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

<<<<<<< HEAD
//GetReader 获取
=======
>>>>>>> e207abb0f14858dd24a18e962477dca1545b8cec
func GetReader(content string, charset string) io.Reader {
	if strings.EqualFold(charset, "utf-8") {
		return strings.NewReader(content)
	}
	return transform.NewReader(bytes.NewReader([]byte(content)), simplifiedchinese.GBK.NewDecoder())
}

<<<<<<< HEAD
//Convert []byte转换为字符串
func Convert(data []byte, encoding string) (content string, err error) {
	if !strings.EqualFold(encoding, "gbk") && strings.EqualFold(encoding, "gb2312") &&
		!strings.EqualFold(encoding, "utf-8") {
		err = fmt.Errorf("不支持的编码方式：%s", encoding)
		return
	}
	//转换utf-8格式
	if strings.EqualFold(encoding, "utf-8") {
		content = string(data)
		return
	}

	//转换gbk gb2312格式
	buffer, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return
	}
	content = string(buffer)
	return
=======
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
>>>>>>> e207abb0f14858dd24a18e962477dca1545b8cec
}
