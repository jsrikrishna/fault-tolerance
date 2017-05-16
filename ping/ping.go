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
}

func (scheduler *Scheduler) GetBackend() (string, error){
	var available []string
	available = scheduler.AvailableServers
	numberOfServers := len(available)
	if numberOfServers == 0 {
		return "", errors.New("All servers are down, no servers to connect")
	}
	serverNumber := rand.Intn(numberOfServers)
	fmt.Printf("Server Number is %d\n", serverNumber)
	return available[serverNumber], nil
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
		defer conn.Close()
		connected = append(connected, value.Address)
	}
	scheduler.AvailableServers = connected
	scheduler.UnavailableServers = disconnected
	scheduler.DeadServers = dead
}
