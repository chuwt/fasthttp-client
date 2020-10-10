package fasthttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	defaultTimeDuration = time.Second * 30 // 默认过期时间

	// 默认的content-type	form 表单
	defaultContentType = "application/x-www-form-urlencoded"
	// json格式的body
	jsonContentType = "application/json"
	// 文件上传
	formContentType = "multipart/form-data"

	EmptyUrlErr  = errors.New("empty url")
	EmptyFileErr = errors.New("empty file")

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Client struct {
	proxy   string // set to all requests
	timeout time.Duration
	crt     *tls.Certificate
	opts    *requestOptions
}

func NewClientPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &Client{
				timeout: defaultTimeDuration,
				crt:     nil,
				opts:    newRequestOptions(),
			}
		},
	}
}

func NewClient() *Client {
	return &Client{
		timeout: defaultTimeDuration,
		crt:     nil,
		opts:    newRequestOptions(),
	}
}

func (c *Client) SetProxy(proxy string) *Client {
	c.proxy = proxy
	return c
}

func (c *Client) SetTimeout(duration time.Duration) *Client {
	c.timeout = duration
	return c
}

func (c *Client) SetCrt(certPath, keyPath string) *Client {
	clientCrt, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		// todo handle this err
		clientCrt = tls.Certificate{}
	}
	c.crt = &clientCrt
	return c
}

func (c *Client) AddParam(key, value string) *Client {
	c.opts.params[key] = value
	return c
}

func (c *Client) AddParams(params RequestParams) *Client {
	for key, value := range params {
		c.opts.params[key] = value
	}
	return c
}

func (c *Client) AddHeader(key, value string) *Client {
	c.opts.headers.normal[key] = value
	return c
}

func (c *Client) AddHeaders(headers RequestHeaders) *Client {
	for key, value := range headers {
		c.opts.headers.normal[key] = value
	}
	return c
}

func (c *Client) AddCookie(key, value string) *Client {
	c.opts.headers.cookies[key] = value
	return c
}

func (c *Client) AddCookies(cookies RequestCookies) *Client {
	for key, value := range cookies {
		c.opts.headers.cookies[key] = value
	}
	return c
}

func (c *Client) AddFile(fileName, filePath string) *Client {
	c.opts.files[fileName] = filePath
	return c
}

func (c *Client) AddFiles(files RequestFiles) *Client {
	for key, value := range files {
		c.opts.files[key] = value
	}
	return c
}

func (c *Client) AddBodyStruct(object interface{}) *Client {
	bodyByte, _ := json.Marshal(object)
	c.opts.body = bodyByte
	return c
}

func (c *Client) AddBodyBytes(bodyBytes []byte) *Client {
	c.opts.body = bodyBytes
	return c
}

func (c *Client) Get(url string) (*Response, error) {
	if url == "" {
		return nil, EmptyUrlErr
	}
	params := make([]string, 0)
	for key, value := range c.opts.params {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}
	url = fmt.Sprintf("%s?%s", url, strings.Join(params, "&"))
	return c.call(url, fasthttp.MethodGet, c.opts.headers, nil)
}

func (c *Client) Post(url string) (*Response, error) {
	if url == "" {
		return nil, EmptyUrlErr
	}
	// 需要根据content-type 进行设置
	return c.call(url, fasthttp.MethodPost, c.opts.headers, c.opts.body)
}

func (c *Client) SendFile(url string, options ...RequestOption) (*Response, error) {
	if url == "" {
		return nil, EmptyUrlErr
	}
	if len(c.opts.files) == 0 {
		return nil, EmptyFileErr
	}
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for fileName, filePath := range c.opts.files {
		fileWriter, err := bodyWriter.CreateFormFile(fileName, path.Base(filePath))
		if err != nil {
			return nil, err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		//不要忘记关闭打开的文件
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		_ = file.Close()
	}
	_ = bodyWriter.Close()
	c.opts.headers.normal["content-type"] = bodyWriter.FormDataContentType()

	return c.call(url, fasthttp.MethodPost, c.opts.headers, bodyBuffer.Bytes())
}

func (c *Client) call(url, method string, headers requestHeaders, body []byte) (*Response, error) {
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
		case jsonContentType:
			if body != nil {
				req.SetBody(body)
			}
		default:
			if !strings.Contains(contentType, formContentType) && body != nil {
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
	// Client.DoTimeout 超时后不会断开连接，所以使用readTimeout
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
