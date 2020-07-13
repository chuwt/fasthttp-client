package fasthttp

import "time"

var defaultClient = NewClient()

func Get(url string, options ...RequestOption) (*Response, error) {
	return defaultClient.Get(url, options...)
}

func Post(url string, body interface{}, options ...RequestOption) (*Response, error) {
	return defaultClient.Post(url, body, options...)
}

func SendFile(url string, options ...RequestOption) (*Response, error) {
	return defaultClient.SendFile(url, options...)
}

func SetTimeout(duration time.Duration) error {
	return defaultClient.SetTimeout(duration)
}

func SetCrt(certPath, keyPath string) error {
	return defaultClient.SetCrt(certPath, keyPath)
}
