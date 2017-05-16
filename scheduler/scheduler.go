package scheduler

import (
	"fault-tolerance/config"
	"fault-tolerance/ping"
)

func New(config config.Configuration) *ping.Scheduler {
	scheduler := ping.Scheduler{
		Servers: config.Servers,
		PingInterval: config.PingInterval,
		HealthcheckInterval:config.HealthcheckInterval,
		StatusCounter: config.StatusCounter,
		AvailableServers: []string{},
		UnavailableServers: []string{},
		DeadServers: []string{},
		Algorithm: config.Algorithm,
		PreviousServer: int(0),
		CurrentServerCounter: int(0),
		AvailableServerPtrs: config.Servers,
	}

	for _, value := range scheduler.Servers {
		value.CurrentCounter = 0
		value.Status = false
		value.Dead = false
	}
	return &scheduler
}

