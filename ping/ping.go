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
}

func (scheduler *Scheduler) GetBackend() (string, error){
	var available []string
	available, _ = Healthcheck(scheduler)

	numberOfServers := len(available)
	if numberOfServers == 0 {
		return "", errors.New("All servers are down, no servers to connect")
	}
	serverNumber := rand.Intn(numberOfServers)
	fmt.Printf("Server Number is %d\n", serverNumber)
	return available[serverNumber], nil
}

func Healthcheck(scheduler *Scheduler) ([]string, []string) {

	flagTimeout := 10
	var connected []string
	var disconnected []string
	servers := scheduler.Servers
	for _, value := range servers {

		conn, err := net.DialTimeout("tcp", value.Address, time.Duration(flagTimeout) * time.Second)
		if err != nil {
			fmt.Printf(value.Address + " Disconnected\n")
			value.Status = false
			disconnected = append(disconnected, value.Address)
			continue
		}
		fmt.Printf(value.Address + " Connected\n")
		value.Status = true
		defer conn.Close()
		connected = append(connected, value.Address)
	}
	return connected, disconnected
}
