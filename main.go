// helpful resource with additional content:
// http://stackoverflow.com/questions/21270945/how-to-read-the-response-from-a-newsinglehostreverseproxy
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/johnlonganecker/go-proxy/proxy"
)

func main() {

	// handle flags
	listeningPort := flag.String("listening-port", "8080", "a string")
	destPort := flag.String("dest-port", "8080", "a string")
	destUrl := flag.String("dest-url", "127.0.0.1", "a string")

	flag.Parse()

	proxy := proxy.Proxy{
		DestPort: *destPort,
		DestUrl:  *destUrl,
	}

	log.Println("Al Broker")
	log.Println("------------------")
	log.Println("listening port: " + *listeningPort)
	log.Println("destination port: " + *destPort)
	log.Println("destination url: " + *destUrl)

	// intercept call
	http.HandleFunc("/v2/catalog", handleCatalog)

	// all other traffic pass on
	http.HandleFunc("/", proxy.HttpHandler)
	http.ListenAndServe(":"+*listeningPort, nil)
}

func handleCatalog(w http.ResponseWriter, r *http.Request) {

	headers := make(map[string]string)

	headers["Content-Type"] = "application/json"

	resp, err := MakeRequest("GET", "http://localhost:8000/json.json", headers, "")
	dec := json.NewDecoder(resp.Body)

	defer resp.Body.Close()

	for {
		var dataMap map[string]interface{}
		if err := dec.Decode(&dataMap); err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%+v", dataMap)
	}

	if err != nil {
		// handle error
	}
}

func MakeRequest(method string, uri string, headers map[string]string, fields string) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, nil)

	// built up request headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// initiate the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
