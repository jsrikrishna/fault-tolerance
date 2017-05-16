package ping

import (
	"fmt"
	"net"
	"time"
	"errors"
	"math/rand"
	"fault-tolerance/config"
)


type Scheduler struct {
	Servers []config.Server
	PingInterval int
	HealthcheckInterval int
	StatusCounter int
	AvailableServers []string
	UnavailableServers []string
	DeadServers []string
	Algorithm string
	PreviousServer int
	CurrentServerCounter int
	AvailableServerPtrs []config.Server
}

func (scheduler *Scheduler) GetBackend() (string, error){
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
		serverNumber := ((scheduler.PreviousServer+1) % numberOfServers)
		scheduler.PreviousServer = serverNumber
		fmt.Printf("Server Number is %d\n", serverNumber)
		return available[serverNumber], nil

	case "weightedroundrobin":
		fmt.Printf("Using Weighted Round Robin Algorithm\n")
		temp := scheduler.PreviousServer
		if (scheduler.CurrentServerCounter > scheduler.AvailableServerPtrs[scheduler.PreviousServer].Weight){
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
	for{
		Healthcheck(scheduler)
		time.Sleep(time.Duration(scheduler.HealthcheckInterval) *time.Millisecond)
	}
}

func Healthcheck(scheduler *Scheduler) () {

	var connected []string
	var disconnected []string
	var dead []string
	var available_servers []config.Server
	servers := scheduler.Servers
	pingInterval := scheduler.PingInterval
	for _, value := range servers {

		conn, err := net.DialTimeout("tcp", value.Address, time.Duration(pingInterval) * time.Second)
		if err != nil {
			fmt.Printf(value.Address + " Disconnected\n")
			value.CurrentCounter+=1
			if value.CurrentCounter >= scheduler.StatusCounter {
				value.Dead = true
				dead = append(dead,value.Address)
			}
			value.Status = false
			disconnected = append(disconnected, value.Address)
			continue
		}
		fmt.Printf(value.Address + " Connected\n")
		value.Status = true
		value.CurrentCounter = 0
		value.Dead = false
		available_servers = append(available_servers,value)
		defer conn.Close()
		connected = append(connected, value.Address)
	}
	scheduler.AvailableServers = connected
	scheduler.UnavailableServers = disconnected
	scheduler.DeadServers = dead
	scheduler.AvailableServerPtrs = available_servers
}
