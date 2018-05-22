package pipeline

import (
	"github.com/pkg/errors"
)

type proxyPipe struct {
	getter chan ProxyItem
}

type ProxyItem struct {
	Addr     string
	User     string
	Password string
}

func (p *proxyPipe) New(items []ProxyItem) *proxyPipe {

	proxy := &proxyPipe{
		getter: make(chan ProxyItem),
	}

	for _, proxyItem := range items {
		proxy.getter <- proxyItem
	}

	return proxy
}

func (p *proxyPipe) get() (proxyItem ProxyItem, err error) {

	proxyItem, ok := <-p.getter

	if !ok {
		return proxyItem, errors.New("can't get proxyPipe item")
	}

	p.getter <- proxyItem
	return proxyItem, err
}
