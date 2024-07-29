package proxy

import "sync"

type ProxyInterface interface {
	Get() string
}

type proxy struct {
	proxies      []string
	counter      int
	counterMutex sync.Mutex
}

func NewProxy(
	proxies []string,
) ProxyInterface {
	return &proxy{
		proxies:      proxies,
		counter:      0,
		counterMutex: sync.Mutex{},
	}
}

func (p *proxy) Get() string {
	var response string

	p.counterMutex.Lock()

	response = p.proxies[p.counter]

	p.counter += 1

	if p.counter == len(p.proxies) {
		p.counter = 0
	}

	p.counterMutex.Unlock()

	return response
}
