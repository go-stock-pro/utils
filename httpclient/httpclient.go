package httpclient

import (
	"net/http"
	"net/url"

	"github.com/rohanraj7316/logger"
)

type HttpClient struct {
	client            *http.Client
	reqResLogging     bool
	reqResBodyLogging bool
}

func New(config ...Config) (*HttpClient, error) {
	cfg := configDefault(config...)

	// setting up proxy
	if cfg.UseProxy {
		pProxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, err
		}

		cfg.Transport.Proxy = http.ProxyURL(pProxyURL)
	}

	err := logger.Configure()
	if err != nil {
		return nil, err
	}

	return &HttpClient{
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: cfg.Transport,
		},
		reqResLogging:     cfg.LogReqResEnable,
		reqResBodyLogging: cfg.LogReqResBodyEnable,
	}, nil
}
