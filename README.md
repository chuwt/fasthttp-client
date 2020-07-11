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
    client = new(Client)
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
    ```

- post

  ```
  client = new(Client)
  resByte, err = client.Post("http://httpbin.org/post",
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
  client = new(Client)
  client.SetTimeout(time.Minute)
  resByte, err = client.SendFile("http://httpbin.org/post",
  		AddFile("a", "/Users/chuwt/Downloads/test.jpg"),
  		AddFile("b", "/Users/chuwt/Downloads/test.jpg"),
  	)
  // AddFile(fileName, filePath)
  ```
- tls

  ```
  client = new(FastHttp)
  client.SetCrt(certPath, certKey)
  ```

- ps
    - client.SetTimeout 非线程安全，倾向于做为全局配置使用
    - client.SetCrt 非线程安全，倾向于全局做为配置使用

# todo
- 返回 body stream 支持
    - 暂时未找到可行办法，看了源码，遇到chunked时是等待所有数据返回后才返回responseBody，所以无法像net/http那样获取到io.Reader
      浪费了我一天的时间😂，不行就要爆改代码了