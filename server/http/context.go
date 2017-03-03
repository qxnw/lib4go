package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fmt"

	"github.com/arsgo/lib4go/utility"
	"github.com/qxnw/lib4go/encoding"
	"github.com/qxnw/lib4go/logger"
)

//Context 上下文
type Context struct {
	StartTime    time.Time
	Response     http.ResponseWriter
	Request      *http.Request
	Session      string
	Address      string
	Script       string
	Encoding     string
	QueryString  map[string]string
	RequestRaw   string
	ResponseCode int
	Log          logger.ILogger
	Method       string
}

//NewContext 构建web请求上下文
func NewContext(loggerName string, w http.ResponseWriter, r *http.Request, address string, script string) (context *Context) {
	context = &Context{Response: w, Request: r, Address: address, Script: script}
	context.StartTime = time.Now()
	context.Method = r.Method
	context.Session = utility.GetGUID()[:8]
	context.Log = logger.GetSession(loggerName, context.Session)
	context.Log.SetTag("script", script)
	return

}

//PassTime 计算当前使用已过去的时间
func (c *Context) PassTime() time.Duration {
	return time.Now().Sub(c.StartTime)
}

//Parse 转换请求参数
func (c *Context) Parse() (err error) {
	c.RequestRaw, c.QueryString, err = c.parse()
	return
}

//Write 写入请求
func (c *Context) Write(code int, content string) {
	c.ResponseCode = code
	c.Response.WriteHeader(code)
	c.Response.Write([]byte(content))
}

//Redirect 页面转跳
func (c *Context) Redirect(code int, url string) {
	c.Response.Header().Set("Location", url)
	c.Response.WriteHeader(code)
	c.Response.Write([]byte("Redirecting to: " + url))
}

func (c *Context) parse() (body string, rt map[string]string, err error) {
	body, err = c.getBodyText()
	if err != nil {
		err = fmt.Errorf("获取请求参数:%v", err)
		return
	}
	c.Request.ParseForm()
	params, err := c.getPostValues(body)
	if err != nil {
		err = fmt.Errorf("解析获取请求参数:%v", err)
		return
	}
	if len(c.Request.Form) > 0 {
		for k, v := range c.Request.Form {
			if len(v) > 0 && len(v[0]) > 0 && !strings.EqualFold(v[0], "") {
				params[k], err = encoding.Convert([]byte(v[0]), c.Encoding)
				if err != nil {
					err = fmt.Errorf("编码转换出错:%v", err)
					return
				}
			}
		}
	}
	return body, params, err
}

func (c *Context) getBodyText() (content string, err error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	content, err = encoding.Convert(body, c.Encoding)
	return
}
func (c *Context) getPostValues(body string) (rt map[string]string, err error) {
	rt = make(map[string]string)
	values, err := url.ParseQuery(body)
	if err != nil {
		return
	}
	for i, v := range values {
		if len(v) > 0 && !strings.EqualFold(v[0], "") {
			rt[i], err = encoding.Convert([]byte(v[0]), c.Encoding)
			if err != nil {
				return
			}
		}
	}
	return
}
