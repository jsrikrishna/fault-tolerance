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
	. "fault-tolerance/routes"
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
		// Reference : https://stackoverflow.com/questions/23070876/reading-body-of-http-request
		// -without-modifying-request-state
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
	}

}

func main() {
	configuration, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		router := NewRouter()
		fmt.Printf("Configuration - %v\n", configuration)
		tracker := NewRequestTracker()
		loadScheduler := scheduler.New(configuration, tracker)
		// Run a goroutine for healthcheck
		go ping.HealthCheckWrapper(loadScheduler)
		proxy := NewMultipleHostReverseProxy(loadScheduler, tracker)
		log.Fatal(http.ListenAndServe(configuration.BindTo, proxy))
		//go func(){log.Fatal(http.ListenAndServe(configuration.BindTo, proxy))}()
		log.Fatal(http.ListenAndServe(configuration.BindToStatusServer, router))
	}

}
