package fasthttp

import (
	"testing"
)

var (
	res *Response
	err error
)

func TestGet(t *testing.T) {
	res, err = NewClient().Get("http://httpbin.org/get")
	t.Log(string(res.Body), err)

	res, err = NewClient().
		AddParam("param1", "param1").
		AddParams(RequestParams{
			"param2": "param2",
			"param3": "param3",
		}).
		AddHeader("header1", "header1").
		AddHeaders(RequestHeaders{
			"header2": "header2",
			"header3": "header3",
		}).
		AddCookie("cookie1", "cookie1").
		AddCookies(RequestCookies{
			"cookie2": "cookie2",
			"cookie3": "cookie3",
		}).
		Get("http://httpbin.org/get")
	if err == nil {
		t.Log(string(res.Body))
	}
	t.Log(err)
}

func TestPost(t *testing.T) {
	res, err = NewClient().Post("http://httpbin.org/post")
	t.Log(string(res.Body), err)

	res, err = NewClient().
		AddBodyStruct(
			struct {
				Request string `json:"request" form:"request"`
				Num     int    `json:"num"`
			}{
				Request: "test",
				Num:     1,
			},
		).
		AddHeader("header1", "header1").
		AddHeaders(RequestHeaders{
			"header2": "header2",
			"header3": "header3",
			//"content-type": "application/json",
		}).
		AddCookie("cookie1", "cookie1").
		AddCookies(RequestCookies{
			"cookie2": "cookie2",
			"cookie3": "cookie3",
		}).
		Post("http://httpbin.org/post")

	if err == nil {
		t.Log(string(res.Body))
	}
	t.Log(err)
}

func TestSendFile(t *testing.T) {

	res, err = NewClient().
		AddFile("a", "/Users/chuwt/Downloads/test.jpg").
		AddFiles(
			RequestFiles{
				"b": "/Users/chuwt/Downloads/test.jpg",
			},
		).
		SendFile("http://httpbin.org/post")
	if err == nil {
		t.Log(string(res.Body))
	}
	t.Log(err)
}

func TestSyncPool(t *testing.T) {
	sp := NewClientPool()
	res, err := sp.Get().(*Client).Get("http://httpbin.org/get")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res.Body))
}
