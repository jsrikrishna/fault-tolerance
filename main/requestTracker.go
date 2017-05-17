package main

import (
	"net/url"
	"fmt"
	"strings"
)

type Request struct {
	Path           string
	StartTime      string
	EndTime        string
	BackEndHandler string
}

type RequestTracker struct {
	requestTracker map[string][]Request
}

func NewRequestTracker() *RequestTracker {
	return &RequestTracker{
		requestTracker: make(map[string][]Request),
	}
}

func (tracker *RequestTracker) addRequest(url *url.URL, backend string) {
	if (strings.TrimSpace(url.Path) == "/resources") {
		value, ok := tracker.requestTracker[backend]
		request := Request{
			Path: url.Path,
			StartTime: "StartTime", // Keeping strings just for now, until /requests format is known
			EndTime: "EndTime", // Keeping strings just for now, until /requests format is known
			BackEndHandler: backend,
		}
		if ok {
			value = append(value, request)
			tracker.requestTracker[backend] = value
		} else {
			var requests []Request
			requests = append(requests, request)
			tracker.requestTracker[backend] = requests
		}
		//for key, value := range tracker.requestTracker {
		//	fmt.Println("Key:", key, "Value:", value)
		//}
	}

}