package httpclient

import "net/http"

type Transport struct {
	Proxy http.RoundTripper
}

func (t Transport) logger(req *http.Request) (res *http.Response, err error) {
	res, err = t.Proxy.RoundTrip(req)
	return res, err
}
