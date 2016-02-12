package interceptor

import (
	"encoding/json"
	"net/http"

	"github.com/johnlonganecker/al-broker/config"
	"github.com/johnlonganecker/al-broker/proxy"
)

type Interceptor struct {
	Listen      config.Transaction
	Destination config.Transaction
	Proxy       proxy.Proxy
}

// 1. make initial http request to catalog service
// 2. parse json
// 3. if mapping exist make reverse proxy request to SB
// 4. merge mapping
// 5. write response back to original ResponseWriter
func (i *Interceptor) HandleHttp(w http.ResponseWriter, r *http.Request) {

	resp, err := i.MakeRequest(i.Destination.HttpMethod, i.Destination.Url, i.Destination.Headers)
	if err != nil {
		handleError(w, err)
		return
	}

	results, err := DecodeBody(*resp)
	if err != nil {
		handleError(w, err)
		return
	}

	// if mappings exist we need to talk to SB
	if len(i.Destination.Mappings) > 0 {
		i.CallProxy(w, r, results)
		return
	}

	// write json results back to the response writer
	enc := json.NewEncoder(w)
	err = enc.Encode(results)

	if err != nil {
		handleError(w, err)
		return
	}
}

func (i *Interceptor) CallProxy(w http.ResponseWriter, r *http.Request, serviceCatalog map[string]interface{}) error {

	transport := interceptorTransport{
		Mappings:    i.Destination.Mappings,
		ExtraFields: i.Destination.ExtraFields,
		MapToMerge:  serviceCatalog,
	}

	i.Proxy.Transport = &transport
	i.Proxy.HttpHandler(w, r)

	return nil
}

func (i *Interceptor) MakeRequest(httpMethod string, uri string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, uri, nil)

	req.Header = header

	// initiate the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return resp, err
	}

	return resp, nil
}
