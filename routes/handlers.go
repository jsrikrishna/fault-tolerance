package routes

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
	. "fault-tolerance/requestTracker"
	. "fault-tolerance/ping"
)

type LoadBalancer struct {
	Tracker   *RequestTracker
	Scheduler *Scheduler
}
type statusBody struct {
	BackEnd   string `json:"backend"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

type newServer struct {
	serverName string `json:"serverName"`
	address    string `json:"address"`
	weight     int `json:"weight"`
}

func (loadBalancer *LoadBalancer) Resources(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got the request")
}

func (loadBalancer *LoadBalancer) RequestStatusHandler(w http.ResponseWriter, req *http.Request) {

	if req.Body != nil {
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		bodyForStatus := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		decoder := json.NewDecoder(bodyForStatus)
		var statusData statusBody
		err := decoder.Decode(&statusData)

		if err != nil {
			fmt.Println("Error occurred while decoding the /status body ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		loadBalancer.Tracker.RemoveRequest(statusData.BackEnd, statusData.StartTime, statusData.EndTime)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func (loadBalancer *LoadBalancer) AddServer(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(req.Body)
		bodyForStatus := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		decoder := json.NewDecoder(bodyForStatus)
		var newServer newServer
		err := decoder.Decode(&newServer)

		if err != nil {
			fmt.Println("Error occurred while decoding the /server body ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}
