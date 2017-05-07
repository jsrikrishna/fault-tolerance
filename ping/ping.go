package main

import (
	"fmt"
	"fault-tolerance/config"
	"net"
	"time"
)


func main() {

	flagTimeout := 10
	var connected []string
	var disconnected []string
	configuration, _ := config.ReadConfig()
	for _, value := range configuration.Servers {

				conn, err := net.DialTimeout("tcp", value.Address, time.Duration(flagTimeout) * time.Second)
				if err != nil {
					fmt.Printf(value.Address + " Disconnected\n")
					disconnected = append(disconnected, value.Address)
					continue
				}
				fmt.Printf(value.Address+ " Connected\n")
				defer conn.Close()
				connected = append(connected,value.Address)
		}
	}
