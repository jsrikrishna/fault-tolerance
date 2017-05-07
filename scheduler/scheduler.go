package scheduler

import "fault-tolerance/config"
import "errors"
import (
	"math/rand"
	"fmt"
)

type Scheduler struct {
	Servers []config.Server
}

func New(config config.Configuration) *Scheduler {
	scheduler := Scheduler{config.Servers}
	return &scheduler
}

func (scheduler *Scheduler) GetBackend() (string, error){
	numberOfServers := len(scheduler.Servers)
	if numberOfServers == 0 {
		return "", errors.New("All servers are down, no servers to connect")
	}
	serverNumber := rand.Intn(numberOfServers)
	fmt.Printf("Server Number is %d\n", serverNumber)
	return scheduler.Servers[serverNumber].Address, nil
}

