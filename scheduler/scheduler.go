package scheduler

import (
	"fault-tolerance/config"
	"fault-tolerance/ping"
)

func New(config config.Configuration) *ping.Scheduler {
	scheduler := ping.Scheduler{config.Servers, config.PingInterval,config.HealthcheckInterval, config.StatusCounter,[]string{},[]string{},[]string{}}

	for _,value := range scheduler.Servers{
		value.CurrentCounter = 0
		value.Status = false
		value.Dead = false
	}
	return &scheduler
}

