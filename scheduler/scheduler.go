package scheduler

import (
	"fault-tolerance/config"
	"fault-tolerance/ping"
)

func New(config config.Configuration) *ping.Scheduler {
	scheduler := ping.Scheduler{config.Servers, config.PingInterval,config.HealthcheckInterval,[]string{}}
	return &scheduler
}

