package httpclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

var (
	REQUEST_TIMEOUT = "2s"
)

type Config struct {
	// Timeout gives you timeout for request
	// Default: 30s
	Timeout time.Duration

	// bool flag which help us in configuring proxy
	// Default: false
	UseProxy bool

	// url need to do the proxy
	// Default: nil
	ProxyURL string

	// LogReqResEnable helps in logging request & responses.
	// Default true
	LogReqResEnable bool

	// LogReqResBodyEnable helps in logging request and responses body
	// Default true
	LogReqResBodyEnable bool

	Transport *http.Transport
}

var ConfigDefault = Config{
	UseProxy:            false,
	LogReqResEnable:     true,
	LogReqResBodyEnable: true,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		timeout, err := time.ParseDuration(REQUEST_TIMEOUT)
		if err != nil {
			ConfigDefault.Timeout = 2 * time.Second
		} else {
			ConfigDefault.Timeout = timeout
		}

		return ConfigDefault
	}

	cfg := config[0]
	if cfg.Timeout != 0 {
		ConfigDefault.Timeout = cfg.Timeout
	}

	if !cfg.LogReqResBodyEnable {
		ConfigDefault.LogReqResBodyEnable = false
	}

	if !cfg.LogReqResEnable {
		ConfigDefault.LogReqResEnable = false
	}

	if cfg.Transport != nil {
		ConfigDefault.Transport = cfg.Transport
	}

	return ConfigDefault
}
