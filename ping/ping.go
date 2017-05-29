package ping

import (
	"fmt"
	"net"
	"time"
	"errors"
	"math/rand"
	"fault-tolerance/config"
	. "fault-tolerance/requestTracker"
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
		time.Sleep(time.Duration(scheduler.HealthcheckInterval) * time.Millisecond)
	}
}

func Healthcheck(scheduler *Scheduler) () {

	var connected []string
	var disconnected []string
	var dead []string
	var available_servers []*config.Server
	pingInterval := scheduler.PingInterval
	for _, value := range scheduler.Servers {

		conn, err := net.DialTimeout("tcp", value.Address, time.Duration(pingInterval) * time.Second)
		if err != nil {
			fmt.Printf(value.Address + " Disconnected\n")
			value.CurrentCounter += 1
			fmt.Printf("Current Counter %d\n", value.CurrentCounter)
			if value.CurrentCounter >= scheduler.StatusCounter {
				value.Dead = true
				backEnd, err := scheduler.GetBackend()
				if err != nil {
					fmt.Printf("Could not a get a backend when server is down %v\n", err)
				} else {
					fmt.Printf("Any request for %q will be routed to %q\n", value.Address, backEnd)
					scheduler.RequestTracker.CheckForDeadServerRequests(value.Address, backEnd)
				}
				// Currently ignore, if we don't get a backend
				dead = append(dead, value.Address)
			}
			value.Status = false
			disconnected = append(disconnected, value.Address)
		} else {
			fmt.Printf(value.Address + " Connected\n")
			value.Status = true
			value.CurrentCounter = 0
			value.Dead = false
			available_servers = append(available_servers, value)
			defer conn.Close()
			connected = append(connected, value.Address)
		}

	}
	scheduler.AvailableServers = connected
	scheduler.UnavailableServers = disconnected
	scheduler.DeadServers = dead
	scheduler.AvailableServerPtrs = available_servers
}
