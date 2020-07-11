package fasthttp

import (
	"bufio"
	"testing"
	"time"
)

var (
	client  = Client{}
	resByte []byte
	err     error
	b       *bufio.Reader
)

func TestGet(t *testing.T) {
	resByte, err = client.Get("http://httpbin.org/get",
		AddParam("param1", "param1"),
		AddParams(RequestParams{
			"param2": "param2",
			"param3": "param3",
		}),
		AddHeader("header1", "header1"),
		AddHeaders(RequestHeaders{
			"header2": "header2",
			"header3": "header3",
		}),
		AddCookie("cookie1", "cookie1"),
		AddCookies(RequestCookie{
			"cookie2": "cookie2",
			"cookie3": "cookie3",
		}),
	)
	t.Log(string(resByte), err)
}

func TestPost(t *testing.T) {
	resByte, err = client.Post("http://httpbin.org/post",
		struct {
			Request string `json:"request" form:"request"`
			Num     int    `json:"num"`
		}{
			Request: "test",
			Num:     1,
		},
		AddHeader("header1", "header1"),
		AddHeaders(RequestHeaders{
			"header2": "header2",
			"header3": "header3",
			//"content-type": "application/json",
		}),
		AddCookie("cookie1", "cookie1"),
		AddCookies(RequestCookie{
			"cookie2": "cookie2",
			"cookie3": "cookie3",
		}),
	)
	t.Log(string(resByte), err)
}

func TestSendFile(t *testing.T) {
	client.SetTimeout(time.Minute)
	resByte, err = client.SendFile("http://httpbin.org/post",
		AddFile("a", "/Users/chuwt/Downloads/test.jpg"),
		AddFiles(
			RequestFiles{
				"b": "/Users/chuwt/Downloads/test.jpg",
			},
		),
	)
	t.Log(string(resByte), err)
}
