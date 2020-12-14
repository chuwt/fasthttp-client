package fasthttp

type (
	Mapper map[string]string

	RequestParams  struct{ Mapper }
	RequestHeaders struct{ Mapper }
	RequestCookies struct{ Mapper }
	RequestFiles   struct{ Mapper }
)

func NewParams() Mapper {
	return make(map[string]string)
}

func NewHeaders() Mapper {
	return make(map[string]string)
}

func NewCookies() Mapper {
	return make(map[string]string)
}

func NewFiles() Mapper {
	return make(map[string]string)
}

func (m Mapper) Get(key string) string {
	value, ok := m[key]
	if ok {
		return value
	}
	return ""
}

func (m Mapper) Set(key, value string) Mapper {
	m[key] = value
	return m
}

func newRequestOptions() *requestOptions {
	return &requestOptions{
		files: RequestFiles{Mapper: NewFiles()},
		headers: requestHeaders{
			normal:  RequestHeaders{Mapper: NewHeaders()},
			cookies: RequestCookies{Mapper: NewCookies()},
		},
		params: RequestParams{Mapper: NewParams()},
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
