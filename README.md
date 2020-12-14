[[_TOC_]]

# fasthttp-client
基于[fasthttp](https://github.com/valyala/fasthttp#installapplication/x-www-form-urlencoded)的http请求客户端，用于快速构建http请求

# 更新
- 2020-12-14
	- 返回添加header
	- 重新设计了请求时入参的结构（param，header，cookie等）
	
# 功能
- get
- post
- sendFile
- 支持使用tls
- 支持proxy(全局)

# 快速开始
- get
    ```
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
    		Get("http://httpbin.org/get")
    ```

- post

  ```
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
  ```

- sendFile

  ```
  res, err = NewClient().
  		AddFile("a", "/Users/chuwt/Downloads/test.jpg").
  		AddFiles(
  			NewFiles().
  				Set("b", "/Users/chuwt/Downloads/test.jpg"),
  		).
  		SendFile("http://httpbin.org/post")
  ```
- tls

  ```
  client = NewClient()
  client.SetCrt(certPath, certKey).Get("")
  ```

- ps
    - client.SetTimeout 非线程安全，倾向于做为全局(实例)配置使用
    - client.SetCrt 非线程安全，倾向于做为全局(实例)配置使用

- 说明
    - 暂时不支持get中带有body的请求
    - post 的content-type默认为`application/x-www-form-urlencoded`
    - 根据fasthttp的[issue](https://github.com/valyala/fasthttp/issues/411), client不支持获取返回的类似io.Reader，需要等待所有
    返回都被接收后才返回client.Do, 所以没法支持 `chunked` 返回
    - 不过[这个人](https://github.com/erikdubbelboer)写了一个[demo](https://github.com/erikdubbelboer/fasthttp/commit/69515271036c791b25543da6a4360fadb6b61173)用来支持获取io.Reader的body，但是没有merge到主分支上去
    - [官方回复](https://github.com/valyala/fasthttp/issues/849)

# todo
- 单个请求的proxy 支持