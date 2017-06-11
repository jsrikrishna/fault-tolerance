package routes

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"fault-tolerance/config"
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
	ServerName string `json:"serverName"`
	Address    string `json:"address"`
	Weight     int `json:"weight"`
}

type systemResource struct {
	Address string `json: "address"`
	UsedPercent float64 `json: "usedPercent"`
	Free uint64 `json: "free"`
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
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			fmt.Println("Error occurred while decoding the /server body ", err)
			http.Error(w, err.Error(), 500)
			return
		}
		var serverDetails newServer
		err = json.Unmarshal(b, &serverDetails)
		if err != nil {
			fmt.Println("Error occurred while decoding the /server body ", err)
			http.Error(w, err.Error(), 500)
			return
		}
		defer req.Body.Close()

		fmt.Printf("New server details are %#v\n", serverDetails)
		var configServer config.Server
		configServer.Name = serverDetails.ServerName
		configServer.Address = serverDetails.Address
		configServer.Weight = serverDetails.Weight
		configServer.CurrentCounter = 0
		configServer.Status = true
		configServer.Dead = false
		loadBalancer.Scheduler.Servers = append(loadBalancer.Scheduler.Servers, &configServer)

		w.WriteHeader(http.StatusOK)
		return

	}
}

func (loadBalancer *LoadBalancer) GetSystemResources(w http.ResponseWriter, req *http.Request){
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			fmt.Println("Error occurred while decoding the /systemResources body ", err)
			http.Error(w, err.Error(), 500)
			return
		}
		var systemResource systemResource
		err = json.Unmarshal(b, &systemResource)
		if err != nil {
			fmt.Println("Error occurred while decoding the /server body ", err)
			http.Error(w, err.Error(), 500)
			return
		}
		for i:= 0; i < len(loadBalancer.Scheduler.Servers); i++ {
			if loadBalancer.Scheduler.Servers[i].Address == systemResource.Address {
				loadBalancer.Scheduler.Servers[i].UsedResources = systemResource.UsedPercent
				loadBalancer.Scheduler.Servers[i].FreeResources = systemResource.Free
				break
			}
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}
