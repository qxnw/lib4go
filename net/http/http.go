package http

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/qxnw/lib4go/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//HTTPClient HTTP客户端
type HTTPClient struct {
	client *http.Client
}

//HTTPClientRequest  http请求
type HTTPClientRequest struct {
	headers  map[string]string
	client   *http.Client
	method   string
	url      string
	params   string
	encoding string
}

//NewHTTPClientCert 根据pem证书初始化httpClient
func NewHTTPClientCert(certFile string, keyFile string, caFile string) (client *HTTPClient, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	caData, err := ioutil.ReadFile(caFile)
	if err != nil {
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	ssl := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	ssl.Rand = rand.Reader
	client = &HTTPClient{}
	client.client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   ssl,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 0)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   0,
			ResponseHeaderTimeout: 0,
		},
	}
	return
}

//NewHTTPClient 构建HTTP客户端，用于发送GET POST等请求
func NewHTTPClient() (client *HTTPClient) {
	client = &HTTPClient{}
	client.client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 0)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   0,
			ResponseHeaderTimeout: 0,
		},
	}
	return
}

//NewHTTPClientProxy 根据代理服务器地址创建httpClient
func NewHTTPClientProxy(proxy string) (client *HTTPClient) {
	client = &HTTPClient{}
	client.client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(proxy) //根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
			},
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 0)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   0,
			ResponseHeaderTimeout: 0,
		},
	}
	return
}

//Download 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
//header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *HTTPClient) Download(method string, url string, params string, header map[string]string) (body []byte, status int, err error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(params))
	if err != nil {
		return
	}
	req.Close = true
	for i, v := range header {
		req.Header.Set(i, v)
	}
	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	status = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}

//Save 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
//header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *HTTPClient) Save(method string, url string, params string, header map[string]string, path string) (status int, err error) {
	body, status, err := c.Download(method, url, params, header)
	if err != nil {
		return
	}
	fl, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	defer fl.Close()
	n, err := fl.Write(body)
	if err == nil && n < len(body) {
		err = io.ErrShortWrite
	}
	return
}

//Request 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
//header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *HTTPClient) Request(method string, url string, params string, charset string, header map[string]string) (content string, status int, err error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(params))
	if err != nil {
		return
	}
	req.Close = true
	for i, v := range header {
		req.Header.Set(i, v)
	}
	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	status = resp.StatusCode
	content, err = encoding.Convert(body, charset)
	return
}

//Get http get请求
func (c *HTTPClient) Get(url string, args ...string) (content string, status int, err error) {
	charset := getEncoding(args...)
	resp, err := c.client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	status = resp.StatusCode
	content, err = encoding.Convert(body, charset)
	return
}

//Post http Post请求
func (c *HTTPClient) Post(url string, params string, args ...string) (content string, status int, err error) {
	charset := getEncoding(args...)
	resp, err := c.client.Post(url, "application/x-www-form-urlencoded", encoding.GetReader(params, charset))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	status = resp.StatusCode
	content, err = encoding.Convert(body, charset)
	return
}

func getEncoding(params ...string) (encoding string) {
	if len(params) > 0 {
		encoding = strings.ToUpper(params[0])
		return
	}
	return "UTF-8"
}
func changeEncodingData(encoding string, data []byte) (content string, err error) {
	if !strings.EqualFold(encoding, "GBK") && !strings.EqualFold(encoding, "GB2312") {
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
