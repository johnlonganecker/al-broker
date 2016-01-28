// helpful resource with additional content:
// http://stackoverflow.com/questions/21270945/how-to-read-the-response-from-a-newsinglehostreverseproxy
package main

import (
	"flag"
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

func main() {

	// handle flags
	listeningPort := flag.String("listening-port", "8080", "a string")
	destPort := flag.String("dest-port", "8080", "a string")
	destUrl := flag.String("dest-url", "127.0.0.1", "a string")

	flag.Parse()

	proxy := Proxy{
		DestPort: *destPort,
		DestUrl:  *destUrl,
	}

	fmt.Println("Al Broker")
	fmt.Println("------------------")
	fmt.Println("listening port: " + *listeningPort)
	fmt.Println("destination port: " + *destPort)
	fmt.Println("destination url: " + *destUrl)

	// intercept call
	http.HandleFunc("/sayblah", SayBlah)

	// all other traffic pass on
	http.HandleFunc("/", proxy.HttpHandler)
	http.ListenAndServe(":"+*listeningPort, nil)
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

func SayBlah(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("BLAH"))
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
