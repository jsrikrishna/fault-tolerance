package scheduler

import (
	"fault-tolerance/config"
	"fault-tolerance/ping"
)

func New(config config.Configuration) *ping.Scheduler {
	scheduler := ping.Scheduler{config.Servers, config.PingInterval,config.HealthcheckInterval, config.StatusCounter,[]string{},[]string{},[]string{},config.Algorithm,int(0),int(0),config.Servers}

	for _,value := range scheduler.Servers{
		value.CurrentCounter = 0
		value.Status = false
		value.Dead = false
	}
	return &scheduler
}

