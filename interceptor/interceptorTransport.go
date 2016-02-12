package interceptor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type interceptorTransport struct {
	Mappings    map[string]string
	ExtraFields map[string]string
	MapToMerge  map[string]interface{}
}

func (t *interceptorTransport) MergeMapInto(m map[string]interface{}) map[string]interface{} {
	for key, value := range t.Mappings {
		m[value] = t.MapToMerge[key]
	}
	return m
}

func (t *interceptorTransport) AddFieldsTo(m map[string]interface{}) map[string]interface{} {
	for key, value := range t.ExtraFields {
		m[key] = value
	}
	return m
}

func (t *interceptorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// decode json response
	results, err := DecodeBody(*resp)
	if err != nil {
		return nil, err
	}

	results = t.MergeMapInto(results)
	results = t.AddFieldsTo(results)

	b, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	resp.ContentLength = -1
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))

	return resp, err
}
