package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	Url       *url.URL
	Transport http.RoundTripper
}

// use straight through
func (p *Proxy) HttpHandler(w http.ResponseWriter, r *http.Request) {

	// request url
	fmt.Println(r.URL.EscapedPath())

	proxy := httputil.NewSingleHostReverseProxy(p.Url)

	// did we define a Transport if not use default
	if p.Transport != nil {
		proxy.Transport = p.Transport
	}

	proxy.ServeHTTP(w, r)
}
