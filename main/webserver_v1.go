package main
//
//import (
//	"fmt"
//	"fault-tolerance/config"
//	//"fault-tolerance/scheduler"
//	//"github.com/parnurzeal/gorequest"
//	//"fault-tolerance/ping"
//	"fault-tolerance/scheduler"
//	"fault-tolerance/ping"
//)
//
//func main() {
//	/*
//	Template to make a HTTP Request
//	resp, body, errs := request.Post("http://localhost:8081/compute").End()
//	fmt.Printf("%v, %v, %v", resp, body, errs)
//	 */
//	configuration, err := config.ReadConfig()
//	if err != nil {
//		fmt.Printf("%v\n", err)
//	}
//	done := make(chan string)
//	loadScheduler := scheduler.New(configuration)
//	tcpServer := New(configuration, loadScheduler, done)
//
//	// Run a goroutine for healthcheck
//	go ping.HealthCheckWrapper(tcpServer.scheduler)
//
//	tcpServer.Start()
//	fmt.Printf("Configuration - %v\n", configuration)
//	<-done
//}
//
