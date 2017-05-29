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
		//for key, value := range tracker.requestTracker {
		//	fmt.Println("Key:", key, "Value:", value)
		//}
	}
}
func currentRequestCallback(resp gorequest.Response, body string, errs []error){
	fmt.Println("Status is done ", resp.Status)
}

func (tracker *RequestTracker) CheckForDeadServerRequests(address string, backend string) {
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
				request := gorequest.New()
				resp, body, errs := request.Post("http://" + backend + "/resources").
					Send(requestBody).
					End(currentRequestCallback)
				fmt.Printf("%v, %v, %v", resp, body, errs)
			}
		}

	}

}
