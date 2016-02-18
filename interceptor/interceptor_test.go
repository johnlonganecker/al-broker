package interceptor_test

import (
	"net/http"
	"net/http/httptest"
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
		interceptor          Interceptor
		proxyUrl             *url.URL
		serviceCatalogServer *ghttp.Server
		serviceBrokerServer  *ghttp.Server
		w                    *httptest.ResponseRecorder
		r                    *http.Request
		err                  error
	)

	BeforeEach(func() {
		serviceCatalogServer = ghttp.NewServer()
		serviceBrokerServer = ghttp.NewServer()

		proxyUrl, _ = url.Parse(serviceCatalogServer.URL())

		/*
			serviceCatalogServer.AppendHandlers(
			)
		*/

		interceptor = Interceptor{
			Listen: config.Transaction{
				Url:         "/v2/catalog",
				HttpMethod:  "GET",
				Headers:     http.Header{"Content-Type": []string{`application/json`}},
				Mappings:    make(map[string]string),
				ExtraFields: make(map[string]string),
			},
			Destination: config.Transaction{
				Url:         serviceCatalogServer.URL() + "/v2/catalog",
				HttpMethod:  "GET",
				Headers:     http.Header{},
				Mappings:    make(map[string]string),
				ExtraFields: make(map[string]string),
			},
			Proxy: proxy.Proxy{
				Url: proxyUrl,
			},
		}
	})

	AfterEach(func() {
		serviceCatalogServer.Close()
		serviceBrokerServer.Close()
	})

	Describe("Make an http requests", func() {
		Context("simulate a catalog fetch request", func() {

			BeforeEach(func() {
				serviceCatalogServer.AppendHandlers(
					ghttp.VerifyRequest("GET", "/v2/catalog"),
				)
			})

			It("should return with a status code of 200", func() {
				resp, err := interceptor.MakeRequest("GET", serviceCatalogServer.URL()+"/v2/catalog", http.Header{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				Expect(serviceCatalogServer.ReceivedRequests()).Should(HaveLen(1))
			})
		})
		Context("simulate a create service request", func() {

			BeforeEach(func() {
				serviceCatalogServer.AppendHandlers(
					ghttp.VerifyRequest("PUT", "/v2/service_instances"),
				)
			})

			It("should return a status code of 200", func() {
				resp, err := interceptor.MakeRequest("PUT", serviceCatalogServer.URL()+"/v2/service_instances", http.Header{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
				Expect(serviceCatalogServer.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("Make http request through proxy package", func() {
		Context("simulate a fetch catalog request", func() {

			BeforeEach(func() {
				serviceCatalogServer.AppendHandlers(
					ghttp.VerifyRequest("GET", "/v2/catalog"),
				)
				w = httptest.NewRecorder()
				r, err = http.NewRequest("GET", serviceCatalogServer.URL()+"/v2/catalog", nil)
			})

			It("should return a status code of 200", func() {
				err := interceptor.CallProxy(w, r, make(map[string]interface{}))
				Expect(err).NotTo(HaveOccurred())
				Expect(w.Code).Should(Equal(http.StatusOK))
				Expect(serviceCatalogServer.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("simulate a create service request", func() {

			BeforeEach(func() {
				serviceCatalogServer.AppendHandlers(
					ghttp.VerifyRequest("PUT", "/v2/service_instances"),
				)
				w = httptest.NewRecorder()
				r, err = http.NewRequest("PUT", serviceCatalogServer.URL()+"/v2/service_instances", nil)
			})

			It("should return a status code of 200", func() {
				err := interceptor.CallProxy(w, r, make(map[string]interface{}))
				Expect(err).NotTo(HaveOccurred())
				Expect(w.Code).Should(Equal(http.StatusOK))
				Expect(serviceCatalogServer.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})

	Describe("Simulate incoming request to the http handler", func() {
		Context("simulate a fetch catalog request", func() {

			BeforeEach(func() {
				serviceCatalogServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, `{
  "services": [{
    "id": "service-guid-here",
    "name": "mysql",
    "description": "A MySQL-compatible relational database",
    "bindable": true,
    "plans": [{
      "id": "plan1-guid-here",
      "name": "small",
      "description": "A small shared database with 100mb storage quota and 10 connections"
    },{
      "id": "plan2-guid-here",
      "name": "large",
      "description": "A large dedicated database with 10GB storage quota, 512MB of RAM, and 100 connections",
       "free": false
    }],
    "dashboard_client": {
      "id": "client-id-1",
      "secret": "secret-1",
      "redirect_uri": "https://dashboard.service.com"
    }
  }]
}`),
						ghttp.VerifyRequest("GET", "/v2/catalog"),
					),
				)
				w = httptest.NewRecorder()
				r, err = http.NewRequest("GET", serviceCatalogServer.URL()+"/v2/catalog", nil)
			})

			It("should return a status code of 200", func() {
				interceptor.HandleHttp(w, r)
				Expect(w.Code).Should(Equal(http.StatusOK))
				Expect(serviceCatalogServer.ReceivedRequests()).Should(HaveLen(1))
			})
		})

		Context("simulate a create service request", func() {

			BeforeEach(func() {
				interceptor = Interceptor{
					Listen: config.Transaction{
						Url:         "/v2/service_instances",
						HttpMethod:  "PUT",
						Headers:     http.Header{"Content-Type": []string{`application/json`}},
						Mappings:    make(map[string]string),
						ExtraFields: make(map[string]string),
					},
					Destination: config.Transaction{
						Url:         serviceCatalogServer.URL() + "/v2/service_instances",
						HttpMethod:  "PUT",
						Headers:     http.Header{},
						Mappings:    make(map[string]string),
						ExtraFields: make(map[string]string),
					},
					Proxy: proxy.Proxy{
						Url: proxyUrl,
					},
				}

				serviceCatalogServer.AppendHandlers(
					ghttp.VerifyRequest("PUT", "/v2/service_instances"),
				)
				w = httptest.NewRecorder()
				r, err = http.NewRequest("PUT", serviceCatalogServer.URL()+"/v2/service_instances", nil)
			})

			It("should return a status code of 200", func() {
				interceptor.HandleHttp(w, r)
				Expect(w.Code).Should(Equal(http.StatusOK))
				Expect(serviceCatalogServer.ReceivedRequests()).Should(HaveLen(1))
			})
		})
	})
})
