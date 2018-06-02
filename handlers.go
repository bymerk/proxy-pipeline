package pipeline

import (
	"net/http"
	"io"
	"golang.org/x/net/proxy"
	"fmt"
)

func (pipe *Pipeline) getProxy(host string) (proxyItem ProxyItem, err error) {
	

	if proxyInterface, ok := pipe.proxies.Load(host); ok {
		if p, ok := proxyInterface.(*proxyPipe); ok {
			proxyItem, err = p.get()
			return proxyItem, err
		}
	} else {
		p := proxyPipe{}
		phost := p.New(pipe.proxyList)
		pipe.proxies.Store(host, phost)
		proxyItem, err = phost.get()
	}

	return proxyItem, nil
}

func (pipe *Pipeline) handleTunneling(w http.ResponseWriter, r *http.Request) {

	proxyItem, err := pipe.getProxy(r.Host)

	fmt.Println(proxyItem)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	//requestDump, err := httputil.DumpRequest(r, true)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(requestDump))


	dialer, err := proxy.SOCKS5("tcp", proxyItem.Addr, &proxy.Auth{User: proxyItem.User, Password: proxyItem.Password}, proxy.Direct)
	destConnection, err := dialer.Dial("tcp", r.Host)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConnection, _, err := hijacker.Hijack()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}


	go pipe.transfer(destConnection, clientConnection)
	go pipe.transfer(clientConnection, destConnection)

}

func (pipe *Pipeline) transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}



func (pipe *Pipeline) handleHTTP(w http.ResponseWriter, r *http.Request) {

	proxyItem, err := pipe.getProxy(r.Host)

	dialer, err := proxy.SOCKS5("tcp", proxyItem.Addr, &proxy.Auth{User: proxyItem.User, Password: proxyItem.Password}, proxy.Direct)

	transport := http.Transport{
		Dial: dialer.Dial,
	}

	resp, err := transport.RoundTrip(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

