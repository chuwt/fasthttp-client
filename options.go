package fasthttp

type (
	RequestParams  map[string]string
	RequestHeaders map[string]string
	RequestCookies map[string]string
	RequestFiles   map[string]string
)

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
	body    []byte
	Proxy   string
	files   RequestFiles
	headers requestHeaders
	params  RequestParams
}

type requestHeaders struct {
	normal  RequestHeaders
	cookies RequestCookies
}

type RequestOption struct {
	f func(*requestOptions)
}
