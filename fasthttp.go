package fasthttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
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

	EmptyUrlErr = errors.New("empty url")
)

type FastHttp struct {
	timeout time.Duration
}

func (fh *FastHttp) SetTimeout(duration time.Duration) {
	fh.timeout = duration
}

func (fh *FastHttp) Get(url string, options ...RequestOption) ([]byte, error) {
	fh.check()
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
	return fh.call(url, fasthttp.MethodGet, opts.headers, nil)
}

func (fh *FastHttp) Post(url string, body interface{}, options ...RequestOption) ([]byte, error) {
	fh.check()
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
		return fh.call(url, fasthttp.MethodPost, opts.headers, bodyByte)
	}
}

func (fh *FastHttp) SendFile(url string, options ...RequestOption) ([]byte, error) {
	fh.check()

	if url == "" {
		return nil, EmptyUrlErr
	}
	opts := newRequestOptions()
	for _, op := range options {
		op.f(opts)
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

	return fh.call(url, fasthttp.MethodPost, opts.headers, bodyBuffer.Bytes())
}

func (fh *FastHttp) call(url, method string, headers requestHeaders, body []byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 用完需要释放资源

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

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	client := &fasthttp.Client{}
	client.ReadTimeout = fh.timeout
	// client.DoTimeout 超时后不会断开连接，所以使用readTimeout
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

// config check
// now only check timeout
func (fh *FastHttp) check() {
	if fh.timeout == 0 {
		fh.timeout = defaultTimeDuration
	}
}
