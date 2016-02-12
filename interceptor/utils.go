package interceptor

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func DecodeBody(resp http.Response) (map[string]interface{}, error) {

	results := make(map[string]interface{})

	dec := json.NewDecoder(resp.Body)

	defer resp.Body.Close()

	// read in the json response as a stream
	for {
		var dataMap map[string]interface{}
		if err := dec.Decode(&dataMap); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		log.Printf("%+v", dataMap)

		for k, v := range dataMap {
			results[k] = v
		}
	}

	return results, nil
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
