package main

import (
	"fmt"
	"fault-tolerance/config"
	"net/http"
	"net/http/httputil"
	"log"
	"fault-tolerance/scheduler"
	"fault-tolerance/ping"
	. "fault-tolerance/requestTracker"
	"io/ioutil"
	"bytes"
)

func NewMultipleHostReverseProxy(scheduler *ping.Scheduler, tracker *RequestTracker) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		backEnd, err := scheduler.GetBackend()
		if err != nil {
			fmt.Printf("Could not a get a backend %v\n", err)
		}
		fmt.Println("Sending the request to ", backEnd)
		var bodyBytes []byte
		// Reference : https://stackoverflow.com/questions/23070876/reading-body-of-http-request-without-modifying-request-state
		if req.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(req.Body)
			bodyForTracking := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			// Restore the io.ReadCloser to its original state
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			tracker.AddRequest(req.URL, bodyForTracking, backEnd)
		}
		req.URL.Scheme = "http"
		req.URL.Host = backEnd
	}
	return &httputil.ReverseProxy{
		Director:director,
		// Commenting for now, as GetBackend() method will always return an available backend
		// but it is possible that, the backEnd returned by GetBackend(), can be down, from the last ping
		// and requests came in between two conseccutvie pings
		//Transport: &http.Transport{
		//	Proxy : func(req *http.Request) (*url.URL, error) {
		//		fmt.Println("Calling a backend")
		//		return http.ProxyFromEnvironment(req)
		//	},
		//	Dial : func(network, addr string) (net.Conn, error) {
		//		fmt.Println("Calling Dial")
		//		conn, err := (&net.Dialer{
		//			Timeout: 30 * time.Second,
		//			KeepAlive: 30 * time.Second,
		//		}).Dial(network, addr)
		//		if err != nil {
		//			fmt.Println("Error during Dial: ", err.Error())
		//		}
		//		return conn, err
		//	},
		//	TLSHandshakeTimeout: 10 * time.Second,
		//},
	}

}

func main() {
	configuration, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("Configuration - %v\n", configuration)
		tracker := NewRequestTracker()
		loadScheduler := scheduler.New(configuration, tracker)
		// Run a goroutine for healthcheck
		go ping.HealthCheckWrapper(loadScheduler)
		proxy := NewMultipleHostReverseProxy(loadScheduler, tracker)
		log.Fatal(http.ListenAndServe(configuration.BindTo, proxy))
	}

}
