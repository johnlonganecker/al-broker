package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/johnlonganecker/al-broker/proxy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Reverse Proxy", func() {
	var (
		w      *httptest.ResponseRecorder
		r      *http.Request
		server *ghttp.Server
		prox   proxy.Proxy
		err    error
		url    *url.URL
	)

	Context("make a http request", func() {

		BeforeEach(func() {
			server = ghttp.NewServer()

			w = httptest.NewRecorder()
			r, err = http.NewRequest("GET", server.URL()+"/v2/catalog", nil)

			url, err = url.Parse(server.URL())

			prox = proxy.Proxy{
				Url: url,
			}

			server.AppendHandlers(
				ghttp.VerifyRequest("GET", "/v2/catalog"),
			)
		})

		AfterEach(func() {
			server.Close()
		})

		It("should return a status code of 200", func() {
			prox.HttpHandler(w, r)
			Expect(w.Code).Should(Equal(http.StatusOK))
			Expect(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
})
