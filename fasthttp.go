package fasthttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

var (
	defaultTimeDuration = time.Second * 30 // 默认过期时间

	// 默认的content-type	form 表单
	defaultContentType = "application/x-www-form-urlencoded"
	// json格式的body
	dataContentType = "application/json"
	// 文件上传
	formContentType = "multipart/form-data"

	EmptyUrlErr  = errors.New("empty url")
	EmptyFileErr = errors.New("empty file")
)

type client struct {
	proxy   string // set to all requests
	timeout time.Duration
	crt     *tls.Certificate
}

func NewClient() *client {
	return &client{
		timeout: defaultTimeDuration,
		crt:     nil,
	}
}

// unThread safe, prefer global setting
// set global proxy
func (c *client) SetProxy(proxy string) error {
	c.proxy = proxy
	return nil
}

// unThread safe, prefer global setting
func (c *client) SetTimeout(duration time.Duration) error {
	c.timeout = duration
	return nil
}

// unThread safe, prefer global setting
func (c *client) SetCrt(certPath, keyPath string) error {
	clientCrt, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}
	c.crt = &clientCrt
	return nil
}

func (c *client) Get(url string, options ...RequestOption) (*Response, error) {
	if url == "" {
		return nil, EmptyUrlErr
	}
	opts := newRequestOptions()
	for _, op := range options {
		op.f(opts)
	}
	params := make([]string, 0)
	for key, value := range opts.params {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}
	url = fmt.Sprintf("%s?%s", url, strings.Join(params, "&"))
	return c.call(url, fasthttp.MethodGet, opts.headers, nil)
}

func (c *client) Post(url string, body interface{}, options ...RequestOption) (*Response, error) {
	if url == "" {
		return nil, EmptyUrlErr
	}
	opts := newRequestOptions()
	for _, op := range options {
		op.f(opts)
	}
	// 需要根据content-type 进行设置
	if bodyByte, err := json.Marshal(body); err != nil {
		return nil, err
	} else {
		return c.call(url, fasthttp.MethodPost, opts.headers, bodyByte)
	}
}

func (c *client) SendFile(url string, options ...RequestOption) (*Response, error) {
	if url == "" {
		return nil, EmptyUrlErr
	}
	opts := newRequestOptions()
	for _, op := range options {
		op.f(opts)
	}
	if len(opts.files) == 0 {
		return nil, EmptyFileErr
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for fileName, filePath := range opts.files {
		fileWriter, err := bodyWriter.CreateFormFile(fileName, path.Base(filePath))
		if err != nil {
			return nil, err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		//不要忘记关闭打开的文件
		defer file.Close()
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return nil, err
		}
	}
	bodyWriter.Close()
	opts.headers.normal["content-type"] = bodyWriter.FormDataContentType()

	return c.call(url, fasthttp.MethodPost, opts.headers, bodyBuffer.Bytes())
}

func (c *client) call(url, method string, headers requestHeaders, body []byte) (*Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	// set cookie
	for key, value := range headers.cookies {
		req.Header.SetCookie(key, value)
	}
	// set header
	for key, value := range headers.normal {
		req.Header.Set(key, value)
	}

	// set body by content-type, only for !=get
	if !req.Header.IsGet() {
		contentType := string(req.Header.ContentType())
		switch contentType {
		case dataContentType:
			req.SetBody(body)
		default:
			if !strings.Contains(contentType, formContentType) {
				argsMap := make(map[string]interface{})
				if err := json.Unmarshal(body, &argsMap); err != nil {
					return nil, err
				}
				fastArgs := new(fasthttp.Args)
				for key, value := range argsMap {
					fastArgs.Add(key, fmt.Sprintf("%v", value))
				}
				req.SetBody(fastArgs.QueryString())
			} else {
				req.SetBody(body)
			}
		}
	}

	client := &fasthttp.Client{
		ReadTimeout: c.timeout,
	}
	if c.crt != nil {
		client.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{*c.crt},
		}
	}
	if c.proxy != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(c.proxy)
	}
	// client.DoTimeout 超时后不会断开连接，所以使用readTimeout
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	ret := &Response{
		Cookie:     make(RequestCookies),
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	resp.Header.VisitAllCookie(func(key, value []byte) {
		ret.Cookie[string(key)] = string(value)
	})
	return ret, nil
}

type Response struct {
	StatusCode int
	Body       []byte
	Cookie     RequestCookies
}
