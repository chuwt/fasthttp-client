package fasthttp

type (
	RequestParams  map[string]string
	RequestHeaders map[string]string
	RequestCookie  map[string]string
	RequestFiles   map[string]string
)

// add single param for get request
// AddParam("p", "1") means http://www.example.com?p=1
func AddParam(key, value string) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			opts.params[key] = value
		},
	}
}

// add single header for all request
// AddHeader("p", "1") means header["p"] = "1"
func AddHeader(key, value string) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			opts.headers.normal[key] = value
		},
	}
}

// add single cookie for all request
// AddHeader("p", "1") means header["p"] = "1"
func AddCookie(key, value string) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			opts.headers.cookies[key] = value
		},
	}
}

func AddFile(fileName, filePath string) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			opts.files[fileName] = filePath
		},
	}
}

// add multi params for get request
// RequestParams is a map[string]string
func AddParams(params RequestParams) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			for key, value := range params {
				opts.params[key] = value
			}
		},
	}
}

// multi headers
func AddHeaders(headers RequestHeaders) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			for key, value := range headers {
				opts.headers.normal[key] = value
			}
		},
	}
}

// multi cookies
func AddCookies(cookies RequestCookie) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			for key, value := range cookies {
				opts.headers.cookies[key] = value
			}
		},
	}
}

func AddFiles(files RequestFiles) RequestOption {
	return RequestOption{
		f: func(opts *requestOptions) {
			for key, value := range files {
				opts.files[key] = value
			}
		},
	}
}

func newRequestOptions() *requestOptions {
	return &requestOptions{
		files: make(map[string]string),
		headers: requestHeaders{
			normal:  make(map[string]string),
			cookies: make(map[string]string),
		},
		params: make(map[string]string),
	}
}

type requestOptions struct {
	files   RequestFiles
	headers requestHeaders
	params  RequestParams
}

type requestHeaders struct {
	normal  RequestHeaders
	cookies RequestCookie
}

type RequestOption struct {
	f func(*requestOptions)
}
