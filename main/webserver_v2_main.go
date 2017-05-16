package main

import (
	"fmt"
	"fault-tolerance/config"
	"net/http"
	"net/url"
	"net/http/httputil"
	"log"
	"net"
	"time"
	"fault-tolerance/scheduler"
	"fault-tolerance/ping"
)

func NewMultipleHostReverseProxy(scheduler *ping.Scheduler) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		backEnd, err := scheduler.GetBackend()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println("Sending the request to ", backEnd)
		req.URL.Scheme = "http"
		req.URL.Host = backEnd
	}
	return &httputil.ReverseProxy{
		Director:director,
		Transport: &http.Transport{
			Proxy : func(req *http.Request) (*url.URL, error) {
				fmt.Println("Calling a backend")
				return http.ProxyFromEnvironment(req)
			},
			Dial : func(network, addr string) (net.Conn, error) {
				fmt.Println("Calling Dial")
				conn, err := (&net.Dialer{
					Timeout: 30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial(network, addr)
				if err != nil {
					fmt.Println("Error during Dial: ", err.Error())
				}
				return conn, err
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

}

func main() {
	configuration, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	loadScheduler := scheduler.New(configuration)
	proxy := NewMultipleHostReverseProxy(loadScheduler)
	log.Fatal(http.ListenAndServe(configuration.BindTo, proxy))
}
