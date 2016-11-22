package utility

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/qxnw/lib4go/security/md5"
)

// GetSessionID 生成一个SessionID（8位）
func GetSessionID() string {
	id := GetGUID()
	return id[:8]
}

// GetGUID 生成Guid字串
func GetGUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return md5.Encrypt(base64.URLEncoding.EncodeToString(b))
}

// GetLocalIPAddress 获取IP地址
func GetLocalIPAddress(masks ...string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	var ipLst []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ipLst = append(ipLst, ipnet.IP.String())
		}
	}
	if len(masks) == 0 && len(ipLst) > 0 {
		return ipLst[0]
	}
	for _, ip := range ipLst {
		for _, m := range masks {
			if strings.HasPrefix(ip, m) {
				return ip
			}
		}
	}
	return "127.0.0.1"
}

// Escape 把编码 \\u0026，\\u003c，\\u003e 替换为 &,<,>
func Escape(input string) string {
	r := strings.Replace(input, "\\u0026", "&", -1)
	r = strings.Replace(r, "\\u003c", "<", -1)
	r = strings.Replace(r, "\\u003e", ">", -1)
	return r
}

// GetExcPath 第一个参数的路径，结果trim掉后面的参数
func GetExcPath(p ...string) string {
	if len(p) == 0 {
		return ""
	}
	if strings.HasPrefix(p[0], ".") {
		path, err := getExecRoot()
		if err != nil {
			return p[0]
		}
		for i := 1; i < len(p); i++ {
			path = strings.Trim(path, p[i])
		}
		return filepath.Join(path, strings.Trim(p[0], "."))
	}
	return p[0]
}

// // Clone 克隆一个变量
// func Clone(src interface{}) (dst interface{}, err error) {
// 	var buf bytes.Buffer
// 	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
// 		return
// 	}
// 	err = gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
// 	return
// }

// GetMin 获取最大wfhg
func GetMin(d int, x int) int {
	if d > x {
		return x
	}
	return d
}

// GetMax 获取最大wfhg
func GetMax(d int, x int) int {
	if d > x {
		return d
	}
	return x
}

// GetMax2 当d为0时，返回x的值，否则取d,y的最大值
func GetMax2(d int, x int, y int) int {
	if d == 0 {
		return x
	}
	return GetMax(d, y)
}

// CloneMap 克隆一个map
func CloneMap(current map[string]interface{}) map[string]interface{} {
	new := make(map[string]interface{})
	for i, v := range current {
		new[i] = v
	}
	return new
}

// Merge 把iput合并到current，如果没有则添加到current，如果有替换成input里面的
func Merge(current map[string]interface{}, input map[string]interface{}) {
	for i, v := range input {
		current[i] = v
	}
}

// DecodeData 格式化[]byte数据
func DecodeData(encoding string, data []byte) (content string, err error) {
	encoding = strings.ToLower(encoding)
	if !strings.EqualFold(encoding, "gbk") && !strings.EqualFold(encoding, "gb2312") {
		content = string(data)
		return
	}
	buffer, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return
	}
	content = string(buffer)
	return
}
