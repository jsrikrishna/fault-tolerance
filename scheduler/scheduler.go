package scheduler

import (
	"math/rand"
	"fmt"
	"fault-tolerance/config"
	"errors"
	"fault-tolerance/ping"
)

type Scheduler struct {
	Servers []config.Server
}

func New(config config.Configuration) *Scheduler {
	scheduler := Scheduler{config.Servers}
	return &scheduler
}

func (scheduler *Scheduler) GetBackend() (string, error){
	var available []string
	available, _ = ping.Healthcheck(scheduler)

	numberOfServers := len(available)
	if numberOfServers == 0 {
		return "", errors.New("All servers are down, no servers to connect")
	}
	serverNumber := rand.Intn(numberOfServers)
	fmt.Printf("Server Number is %d\n", serverNumber)
	return available[serverNumber], nil
}

