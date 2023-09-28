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

func NewHTTPClientPool() *HTTPClientPool {
	return &HTTPClientPool{clients: []ClientWithProxy{{Client: &http.Client{}, Proxy: ""}}}
}

func NewHTTPClientPoolWithProxies(proxies []string, scheme string) *HTTPClientPool {
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
				Proxy: http.ProxyURL(proxyURL),
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
