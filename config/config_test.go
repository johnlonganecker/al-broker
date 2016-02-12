package config_test

import (
	"net/http"

	. "github.com/johnlonganecker/al-broker/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config File", func() {

	var (
		config                  Config
		fileContents            []byte
		validFilePath           string
		invalidFilePath         string
		invalidYamlFileContents []byte
	)

	BeforeEach(func() {
		validFilePath = "../config.yml"
		invalidFilePath = "blah"

		config = Config{
			Port:   "",
			SbPort: "",
			SbUrl:  "",
			Routes: []Route{
				Route{
					Listen: Transaction{
						Url:         "",
						HttpMethod:  "",
						Headers:     http.Header{},
						Mappings:    make(map[string]string),
						ExtraFields: make(map[string]string),
					},
					Destination: Transaction{
						Url:         "",
						HttpMethod:  "",
						Headers:     http.Header{},
						Mappings:    make(map[string]string),
						ExtraFields: make(map[string]string),
					},
				},
			},
		}

		invalidYamlFileContents = []byte(`field: 
--`)

		fileContents = []byte(`port: 8080
sb_port: 8000
sb_url: http://localhost
routes:
- listen:
    url: /v2/catalog
    http_method: GET
    headers:
      "content-type":
      - "application/json"
  destination:
    url: "http://localhost:8001/catalog-service.json"
    http_method: GET
    headers:
      "content-type":
      - "application/json"
    mappings:
      services: services
    extra_fields:
      plan_updateable: true
- listen:
    url: /v2/service_instances
    http_method: PUT
    headers:
      "content-type":
      - "application/json"
  destination:
    url: "http://localhost:8001/catalog-service.json"
    http_method: GET
    headers:
      "content-type":
      - "application/json"
    mappings:
      services: services`)

	})

	Describe("Load config file", func() {
		Context("with valid file", func() {
			_, err := LoadFile(validFilePath)
			It("should return no error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("with invalid file", func() {
			_, err := LoadFile(invalidFilePath)
			It("should fail if file doesn't exist", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Unmarshalling a config file", func() {
		Context("with a valid yaml file", func() {
			It(" should unmarshal successfully", func() {
				err := Unmarshal(&config, fileContents)
				Expect(err).NotTo(HaveOccurred())
			})
		})
		Context("with an invalid yaml file", func() {
			It("should not unmarshal successfully", func() {
				err := Unmarshal(&config, invalidYamlFileContents)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
