// helpful resource with additional content:
// http://stackoverflow.com/questions/21270945/how-to-read-the-response-from-a-newsinglehostreverseproxy

// one al-broker per SB

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/johnlonganecker/al-broker/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/johnlonganecker/al-broker/config"
	"github.com/johnlonganecker/al-broker/interceptor"
	"github.com/johnlonganecker/al-broker/proxy"
)

func main() {

	// handle flags
	configFile := flag.String("config", "config.yml", "path to config file")

	flag.Parse()

	// load configs
	c := config.Config{}
	err := c.LoadConfigFile(*configFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Service Broker URL
	u, err := url.Parse(c.SbUrl + ":" + c.SbPort)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	proxy := proxy.Proxy{Url: u}

	log.Println("Al Broker")
	log.Println("------------------")
	log.Println("listening port: " + c.Port)
	log.Println("service broker port: " + c.SbPort)
	log.Println("service broker url: " + c.SbUrl)

	muxRouter := mux.NewRouter()

	// listen for each route and apply its handler function
	for _, route := range c.Routes {

		inter := interceptor.Interceptor{}

		inter.Listen = route.Listen
		inter.Destination = route.Destination
		inter.Proxy = proxy

		muxRouter.HandleFunc(route.Listen.Url, inter.HandleHttp).Methods(route.Listen.HttpMethod)

		fmt.Println("adding route " + route.Listen.HttpMethod + " " + route.Listen.Url)
	}

	// all other traffic pass on to
	muxRouter.HandleFunc("/{.*}", proxy.HttpHandler)
	http.ListenAndServe(":"+c.Port, muxRouter)
}
