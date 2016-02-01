package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	DestPort string
	DestUrl  string
}

func (p *Proxy) HttpHandler(w http.ResponseWriter, r *http.Request) {

	// request url
	fmt.Println(r.URL.EscapedPath())

	u, err := url.Parse(p.DestUrl + ":" + p.DestPort)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Transport = &interceptorTransport{}

	proxy.ServeHTTP(w, r)
}

type interceptorTransport struct{}

func (t *interceptorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	log.Print(string(body))

	return resp, err
}
