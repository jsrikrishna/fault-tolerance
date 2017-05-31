package requestTracker

import (
	"time"
	"net/url"
	"io"
	"strings"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

const timeLayout = "Mon, 01/02/06, 03:04PM" // Reference Time Format

type Request struct {
	Path           string
	StartTime      time.Time
	EndTime        time.Time
	//StartTime      string
	//EndTime        string
	BackEndHandler string
}

type RequestTracker struct {
	CurrentRequests map[string][]*Request
}

type RequestType struct {
	Type      string `json:"type"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

func NewRequestTracker() *RequestTracker {
	return &RequestTracker{
		CurrentRequests: make(map[string][]*Request),
	}
}

func (tracker *RequestTracker) AddRequest(reqURL *url.URL, body io.ReadCloser, backend string) {
	if (strings.TrimSpace(reqURL.Path) == "/resources") {
		decoder := json.NewDecoder(body)
		var requestType RequestType
		err := decoder.Decode(&requestType)
		if err != nil {
			fmt.Println("Error occurred while decoding the /resources body ", err)
			return
		}
		startTime, err := time.Parse(timeLayout, requestType.StartTime)
		endTime, err := time.Parse(timeLayout, requestType.EndTime)
		fmt.Println(startTime)
		fmt.Println(endTime)
		if err != nil {
			fmt.Println("Error occurred in formatting time, expected timeFormat is ", timeLayout)
			return
		}
		value, ok := tracker.CurrentRequests[backend]
		request := &Request{
			Path: reqURL.Path,
			StartTime: startTime,
			EndTime: endTime,
			BackEndHandler: backend,
		}
		if ok {
			value = append(value, request)
			tracker.CurrentRequests[backend] = value
		} else {
			var requests []*Request
			requests = append(requests, request)
			tracker.CurrentRequests[backend] = requests
		}
		for key, value := range tracker.CurrentRequests {
			fmt.Printf("Key: %s, Value %+v\n", key, value)
		}
		fmt.Printf("AddRequest Number of elements in map %d\n", len(tracker.CurrentRequests))
	}
}

func (tracker *RequestTracker) RemoveRequest(backend string, startTime string, endTime string) {
	requests, present := tracker.CurrentRequests[backend];
	//for key, value := range tracker.CurrentRequests {
	//	fmt.Printf("Key: %s, Value %+v", key, value)
	//}
	fmt.Printf("RemoveRequest Number of elements in map %d\n", len(tracker.CurrentRequests))
	fmt.Printf("Yes backend %v is present, %v with values %+v\n", backend, present, requests)
	if present {
		startTime, err := time.Parse(timeLayout, startTime)
		endTime, err := time.Parse(timeLayout, endTime)
		if err != nil {
			fmt.Println("Error occurred in formatting time while processing removeRequests, " +
				"expected timeFormat is \n", timeLayout)
			return
		}
		for i, request := range requests {
			fmt.Printf("Given start time is %s and end time is %s\n", startTime.String(), endTime.String())
			fmt.Printf("Map start time is %s and end time is %s\n", request.StartTime.String(), request.EndTime.String())
			if endTime.Equal(request.EndTime) && startTime.Equal(request.StartTime) {
				tracker.CurrentRequests[backend] = append(requests[:i], requests[i + 1:]...)
				fmt.Println("Yes present and removed now")
			}

		}
		for key, value := range tracker.CurrentRequests {
			fmt.Println("Key:", key, "Value:", value)
		}
		if len(tracker.CurrentRequests[backend]) == 0 {
			delete(tracker.CurrentRequests, backend)
		}

	}
}

func currentRequestCallback(resp gorequest.Response, body string, errs []error) {
	fmt.Println("Status is done ", resp.Status)
}

func (tracker *RequestTracker) CheckForDeadServerRequests(address string, otherBackend string) {
	otherBackendValue, ok := tracker.CurrentRequests[otherBackend];
	if !ok {
		var requests []*Request
		otherBackendValue = requests

	}

	if currentRequests, present := tracker.CurrentRequests[address]; present {
		fmt.Printf("Server is currently handling requests %t \n", present)
		fmt.Printf("Current Requets are %d\n", len(currentRequests))
		currentTime := time.Now()
		fmt.Println(currentTime.Format(timeLayout))
		for _, currentRequest := range currentRequests {
			if currentRequest.EndTime.After(currentTime) {
				fmt.Println("Yes, need to process request")
				requestBody := RequestType{
					Type : "resources",
					StartTime: currentTime.Format(timeLayout),
					EndTime: currentRequest.EndTime.Format(timeLayout),
				}
				fmt.Printf("Request Body is %+v\n", requestBody)
				request := gorequest.New()
				resp, body, errs := request.Post("http://" + otherBackend + "/resources").
					Send(requestBody).
					End(currentRequestCallback)
				fmt.Printf("%v, %v, %v\n", resp, body, errs)
			}
			newStartTime, err := time.Parse(timeLayout, currentTime.Format(timeLayout))
			if err != nil {
				fmt.Printf("Time format error occured while storing new requests of dead server %s", address)
			}
			currentRequest.StartTime = newStartTime
			otherBackendValue = append(otherBackendValue, currentRequest)
		}
		tracker.CurrentRequests[otherBackend] = otherBackendValue
		delete(tracker.CurrentRequests, address)
	}

}
