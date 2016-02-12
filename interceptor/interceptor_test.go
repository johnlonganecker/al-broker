package interceptor_test

import (
	"net/http"
	"net/url"

	"github.com/johnlonganecker/al-broker/config"
	. "github.com/johnlonganecker/al-broker/interceptor"
	"github.com/johnlonganecker/al-broker/proxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Interceptor", func() {
	var (
		interceptor Interceptor
		server      *ghttp.Server
		proxy       proxy.Proxy
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		// server.URL()

		interceptor = Interceptor{
			Listen: config.Transaction{
				url:          "/v2/catalog",
				http_method:  "GET",
				headers:      http.Header{"Content-Type", []string{`application/json`}},
				mappings:     make(map[string]string),
				extra_fields: make(map[string]string),
			},
			Destination: config.Transaction{
				url:          "",
				http_method:  "",
				headers:      http.Header{},
				mappings:     make(map[string]string),
				extra_fields: make(map[string]string),
			},
			Proxy: proxy.Proxy{
				Url: url.Parse("http://localhost:8080")
			},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	interceptor.MakeRequest("GET", "www.google.com", http.Header)
})
