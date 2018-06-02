package pipeline

import (
	"sync"
	"net/http"
	"time"
)

type Pipeline struct {
	Server     *http.Server
	proxies    *sync.Map
	proxyList []ProxyItem
}

func New() *Pipeline {
	return &Pipeline{
		proxies: &sync.Map{},
		Server: &http.Server{
			ReadHeaderTimeout: time.Duration(time.Second * 15),
			WriteTimeout:      time.Duration(time.Second * 15),
			Addr:              "127.0.0.1:8888",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodConnect {
					//pipe.handleTunneling(w, r)
				} else {
					//pipe.handleHTTP(w, r)
				}
			}),
		},
	}
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