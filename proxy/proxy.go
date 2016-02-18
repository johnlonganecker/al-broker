package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	Url       *url.URL
	Transport http.RoundTripper
}

func (p *Proxy) HttpHandler(w http.ResponseWriter, r *http.Request) {

	proxy := httputil.NewSingleHostReverseProxy(p.Url)

	// did we define a Transport if not use default
	if p.Transport != nil {
		proxy.Transport = p.Transport
	}

	proxy.ServeHTTP(w, r)
}
