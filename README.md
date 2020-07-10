[[_TOC_]]

# fasthttp-client
基于[fasthttp](https://github.com/valyala/fasthttp#installapplication/x-www-form-urlencoded)的http请求客户端，用于快速构建http请求

# 功能
- get
- post
- sendFile
- 支持使用tls

# 快速开始
- get
    ```
    fh = new(FastHttp)
    resByte, err = fh.Get("http://httpbin.org/get",
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
    ```

- post

  ```
  fh = new(FastHttp)
  resByte, err = fh.Post("http://httpbin.org/post",
  		struct {
  			Request string `json:"request"`
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
  ```

- sendFile

  ```
  fh = new(FastHttp)
  fh.SetTimeout(time.Minute)
  resByte, err = fh.SendFile("http://httpbin.org/post",
  		AddFile("a", "/Users/chuwt/Downloads/test.jpg"),
  		AddFile("b", "/Users/chuwt/Downloads/test.jpg"),
  	)
  // AddFile(fileName, filePath)
  ```
- tls

  ```
  fh = new(FastHttp)
  fh.SetCrt(certPath, certKey)
  ```

- ps
- fh.SetTimeout 非线程安全，倾向于全局配置
- fh.SetCrt 非线程安全，倾向于全局配置