package pipeline

import (
	"sync"
	"net/http"
	"time"
	"gitlab.com/merk/pipe"
)

type Pipeline struct {
	Server     *http.Server
	proxies    *sync.Map
	proxyList []ProxyItem
}

func New() *Pipeline {
	p := &Pipeline{
		proxies: &sync.Map{},
		Server: &http.Server{
			ReadHeaderTimeout: time.Duration(time.Second * 15),
			WriteTimeout:      time.Duration(time.Second * 15),
			Addr:              "127.0.0.1:8888",
		},
	}

	p.Server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			p.handleTunneling(w, r)
		} else {
			p.handleHTTP(w, r)
		}
	})

	return p
}

func (pipe *Pipeline) SetProxyList(proxyList []ProxyItem) {
	pipe.proxyList = proxyList
}

func (pipe *Pipeline) Run() error {
	return pipe.Server.ListenAndServe()
}

func (pipe *Pipeline) RunTLS(cert, key string) error {
	return pipe.Server.ListenAndServeTLS(cert, key)
}