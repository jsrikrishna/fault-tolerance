package scheduler

import (
	. "fault-tolerance/config"
	"fault-tolerance/ping"
	. "fault-tolerance/requestTracker"
)

func New(newConfig Configuration, tracker *RequestTracker) *ping.Scheduler {
	var servers []*Server
	for i := 0; i < len(newConfig.Servers); i += 1 {
		servers = append(servers, &newConfig.Servers[i])
	}
	scheduler := ping.Scheduler{
		Servers: servers,
		PingInterval: newConfig.PingInterval,
		HealthcheckInterval:newConfig.HealthcheckInterval,
		StatusCounter: newConfig.StatusCounter,
		AvailableServers: []string{},
		UnavailableServers: []string{},
		DeadServers: []string{},
		Algorithm: newConfig.Algorithm,
		PreviousServer: int(0),
		CurrentServerCounter: int(0),
		AvailableServerPtrs: servers,
		RequestTracker:tracker,
	}

	for _, value := range scheduler.Servers {
		value.CurrentCounter = 0
		value.Status = false
		value.Dead = false
	}
	return &scheduler
}

