package fasthttp

import (
	_ "net/http/pprof"
	"testing"
)

var (
	res *Response
	err error
)

func TestGet(t *testing.T) {
	res, err = NewClient().Get("http://httpbin.org/get")
	if err == nil {
		t.Log(string(res.Body))
	}
	t.Log(err)

	res, err = NewClient().
		AddParam("param1", "param1").
		AddParams(
			NewParams().
				Set("param2", "param2").
				Set("param3", "param3")).
		AddHeader("header1", "header1").
		AddHeaders(
			NewHeaders().
				Set("header2", "header2").
				Set("header3", "header3")).
		AddCookie("cookie1", "cookie1").
		AddCookies(
			NewCookies().
				Set("cookie1", "cookie1").
				Set("cookie2", "cookie2")).
		Get("http://httpbin.org/get?am=999")
	if err == nil {
		t.Log(string(res.Body))
	}
	t.Log(err)
}

func TestPost(t *testing.T) {
	res, err = NewClient().Post("http://httpbin.org/post")
	if err == nil {
		t.Log(string(res.Body))
	}
	t.Log(err)

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
		AddHeaders(
			NewHeaders().
				Set("header2", "header2").
				Set("header3", "header3"),
			//"content-type": "application/json",
		).
		AddCookie("cookie1", "cookie1").
		AddCookies(NewCookies().
			Set("cookie1", "cookie1").
			Set("cookie2", "cookie2")).
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
			NewFiles().
				Set("b", "/Users/chuwt/Downloads/test.jpg"),
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

func TestMapper(t *testing.T) {
	adder := NewHeaders().
		Set("a", "b").
		Set("c", "d")
	t.Log(adder)
}
