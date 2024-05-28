package core

import (
	"net/http"
	"net/url"
	"sync"
)

type ClientWithProxy struct {
	Client *http.Client
	Proxy  string
}

type HTTPClientPool struct {
	clients []ClientWithProxy
	mutex   sync.Mutex
	index   int
}

func NewHTTPClientPool(config *TransportConfig) *HTTPClientPool {
	return &HTTPClientPool{clients: []ClientWithProxy{{Client: &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   config.DisableKeepAlive,
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		},
	}, Proxy: ""}}}
}

func NewHTTPClientPoolWithProxies(proxies []string, scheme string, config *TransportConfig) *HTTPClientPool {
	if len(proxies) == 0 {
		panic("proxies are empty")
	}

	var clients []ClientWithProxy

	for _, proxy := range proxies {
		proxyURL, err := url.Parse(scheme + "://" + proxy)
		if err != nil {
			continue
		}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyURL(proxyURL),
				DisableKeepAlives:   config.DisableKeepAlive,
				MaxIdleConns:        config.MaxIdleConns,
				MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
				IdleConnTimeout:     config.IdleConnTimeout,
			},
		}
		clients = append(clients, ClientWithProxy{Client: client, Proxy: proxy})
	}

	return &HTTPClientPool{clients: clients}
}

func (h *HTTPClientPool) Get() *ClientWithProxy {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	clientWithProxy := h.clients[h.index]
	h.index = (h.index + 1) % len(h.clients)
	return &clientWithProxy
}
