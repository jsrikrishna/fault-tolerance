package ping

import (
	"fmt"
	"net"
	"time"
	"errors"
	"math/rand"
	"fault-tolerance/config"
	. "fault-tolerance/requestTracker"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"math"
)

type Scheduler struct {
	Servers              []*config.Server
	PingInterval         int
	HealthcheckInterval  int
	StatusCounter        int
	AvailableServers     []string
	UnavailableServers   []string
	DeadServers          []string
	Algorithm            string
	PreviousServer       int
	CurrentServerCounter int
	AvailableServerPtrs  []*config.Server
	RequestTracker       *RequestTracker
}

func (scheduler *Scheduler) GetBackend() (string, error) {
	var available []string
	available = scheduler.AvailableServers
	numberOfServers := len(available)
	if numberOfServers == 0 {
		return "", errors.New("All servers are down, no servers to connect")
	}

	switch scheduler.Algorithm {
	case "random":
		fmt.Printf("Using Random Algorithm\n")
		serverNumber := rand.Intn(numberOfServers)
		fmt.Printf("Server Number is %d\n", serverNumber)
		return available[serverNumber], nil

	case "roundrobin":
		fmt.Printf("Using Round Robin Algorithm\n")
		serverNumber := ((scheduler.PreviousServer + 1) % numberOfServers)
		scheduler.PreviousServer = serverNumber
		fmt.Printf("Server Number is %d\n", serverNumber)
		return available[serverNumber], nil

	case "weightedroundrobin":
		fmt.Printf("Using Weighted Round Robin Algorithm\n")
		temp := scheduler.PreviousServer
		if (scheduler.CurrentServerCounter > scheduler.AvailableServerPtrs[scheduler.PreviousServer].Weight) {
			temp = (scheduler.PreviousServer + 1) % numberOfServers
			scheduler.PreviousServer = (scheduler.PreviousServer + 1) % numberOfServers
			scheduler.CurrentServerCounter = 0

		}
		scheduler.CurrentServerCounter++
		return available[temp], nil

	case "cpumetrics":
		fmt.Printf("Using CPU Metrics")
		min := uint64(math.MaxUint64)
		serverIndex := 0
		for i := 0; i < numberOfServers; i++ {
			if (scheduler.AvailableServerPtrs[i].FreeResources < min) {
				serverIndex = i
				min = scheduler.AvailableServerPtrs[i].FreeResources
			}
		}
		fmt.Printf("Using Server %d\n", serverIndex)
		return available[serverIndex], nil


	default:
		fmt.Printf("Defaulting to Random Algorithm\n")
		serverNumber := rand.Intn(numberOfServers)
		fmt.Printf("Server Number is %d\n", serverNumber)
		return available[serverNumber], nil
	}

}

func HealthCheckWrapper(scheduler *Scheduler) {
	for {
		Healthcheck(scheduler)
		SystemResourceChecks(scheduler)
		time.Sleep(time.Duration(scheduler.HealthcheckInterval) * time.Millisecond)
	}
}

func Healthcheck(scheduler *Scheduler) () {

	var connected []string
	var disconnected []string
	var dead []string
	var available_servers []*config.Server
	pingInterval := scheduler.PingInterval
	var deleted = 0

	for index, _ := range scheduler.Servers {
		j := index - deleted
		curr_server := scheduler.Servers[j]
		if (curr_server.Dead) {
			continue
		}

		conn, err := net.DialTimeout("tcp", curr_server.Address, time.Duration(pingInterval) * time.Second)
		if err != nil {
			fmt.Printf(curr_server.Address + " Disconnected\n")
			curr_server.CurrentCounter += 1
			fmt.Printf("Current Counter %d\n", curr_server.CurrentCounter)
			if curr_server.CurrentCounter >= scheduler.StatusCounter {
				curr_server.Dead = true
				// sleep for HealthCheckInterval time + 1sec
				time.Sleep(time.Duration(scheduler.HealthcheckInterval + 1000) * time.Millisecond)
				backEnd, err := scheduler.GetBackend()
				if err != nil {
					fmt.Printf("Could not a get a backend when server is down %v\n", err)
				} else {
					fmt.Printf("Any request for %q will be routed to %q\n", curr_server.Address, backEnd)
					scheduler.RequestTracker.CheckForDeadServerRequests(curr_server.Address, backEnd)
				}
				// Currently ignore, if we don't get a backend
				dead = append(dead, curr_server.Address)
				scheduler.Servers = scheduler.Servers[:j + copy(scheduler.Servers[j:], scheduler.Servers[j + 1:])]
				deleted++
			}
			curr_server.Status = false
			disconnected = append(disconnected, curr_server.Address)
		} else {
			fmt.Printf(curr_server.Address + " Connected\n")
			curr_server.Status = true
			curr_server.CurrentCounter = 0
			curr_server.Dead = false
			available_servers = append(available_servers, curr_server)
			defer conn.Close()
			connected = append(connected, curr_server.Address)
		}

	}
	scheduler.AvailableServers = connected
	scheduler.UnavailableServers = disconnected
	scheduler.DeadServers = dead
	scheduler.AvailableServerPtrs = available_servers
}

func SystemResourceChecks(scheduler *Scheduler) {

	for i := 0; i < len(scheduler.Servers); i++ {
		serverAddress := scheduler.Servers[i].Address
		request := gorequest.New()
		resp, body, errs := request.Get("http://" + serverAddress + "/systemResources").End()
		if errs != nil {
			continue
		}
		if resp.StatusCode == 200 {
			var systemResourcesData SystemResourceData
			if err := json.Unmarshal([]byte(body), &systemResourcesData); err != nil {
				fmt.Printf("Error occurred when decoding system resource %#v", err)
				continue
			}
			scheduler.Servers[i].UsedResources = systemResourcesData.UsedPercent
			scheduler.Servers[i].FreeResources = systemResourcesData.Free
			fmt.Printf("Res type is %#v\n", systemResourcesData)
		} else {
			continue
		}

	}
}
