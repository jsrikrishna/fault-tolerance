package main

import (
	"strings"
	//"net/http"
	"encoding/json"
	"fmt"
	"time"
	"net/url"
	"io"
)

type Request struct {
	Path           string
	StartTime      time.Time
	EndTime        time.Time
	//StartTime      string
	//EndTime        string
	BackEndHandler string
}

type RequestTracker struct {
	requestTracker map[string][]Request
}

type RequestType struct {
	Type      string `json:"type"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

func NewRequestTracker() *RequestTracker {
	return &RequestTracker{
		requestTracker: make(map[string][]Request),
	}
}

func (tracker *RequestTracker) addRequest(reqURL *url.URL, body io.ReadCloser, backend string) {
	if (strings.TrimSpace(reqURL.Path) == "/resources") {
		decoder := json.NewDecoder(body)
		var requestType RequestType
		err := decoder.Decode(&requestType)
		if err != nil {
			fmt.Println("Error occurred while decoding the /resources body ", err)
			return
		}
		timeLayout := "Mon, 01/02/06, 03:04PM" // Reference Time Format
		startTime, err := time.Parse(timeLayout, requestType.StartTime)
		endTime, err := time.Parse(timeLayout, requestType.EndTime)
		fmt.Println(startTime)
		fmt.Println(endTime)
		if err != nil {
			fmt.Println("Error occurred in formatting time, expected timeFormat is ", timeLayout)
			return
		}
		value, ok := tracker.requestTracker[backend]
		request := Request{
			Path: reqURL.Path,
			StartTime: startTime,
			EndTime: endTime,
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