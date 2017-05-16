package main

import (
	"net/url"
	"fmt"
)

type Request struct {
	Path           string
	StartTime      string
	EndTime        string
	BackEndHandler string
}

type RequestTracker struct {
	requestTracker map[string]Request
}

func NewRequestTracker() *RequestTracker {
	return &RequestTracker{
		requestTracker: make(map[string]Request),
	}
}

func (*RequestTracker) addRequest(url *url.URL, backend string) {
	fmt.Println("Request went to", backend, " with path ", url.Path)
}