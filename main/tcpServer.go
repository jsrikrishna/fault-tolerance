package main

import (
	"net"
	"fmt"
	"fault-tolerance/config"
	"time"
	"fault-tolerance/scheduler"
)

type TcpServer struct {
	name       string
	bindTo     string
	listener   net.TCPListener
	connect    chan (*TcpContext)
	disconnect chan (net.Conn)
	done       chan string
	scheduler *scheduler.Scheduler
}

func New(configuration config.Configuration, scheduler *scheduler.Scheduler, done chan string) *TcpServer {
	tcpServer := &TcpServer{
		name: configuration.Name,
		bindTo:configuration.BindTo,
		done: done,
		connect: make(chan *TcpContext),
		disconnect:make(chan net.Conn),
		scheduler:scheduler,
	}
	return tcpServer
}
func (server *TcpServer) Start() (err error) {
	go func() {
		for {
			select {
			case client := <-server.disconnect:
				fmt.Println("Coming here now")
				server.HandleClientDisConnect(client)
			case context := <-server.connect:
				server.HandleClientConnect(context)

			}
		}

	}()

	if err := server.Listen(server.bindTo, server.done); err != nil {
		fmt.Printf("There is an error while listening %v\n", err)
		return err
	}
	return nil
}

func (server *TcpServer) Listen(bindTo string, done chan string) (err error) {
	listener, err := net.Listen("tcp", bindTo)
	if err != nil {
		fmt.Printf("Error 1 %v", err)
		return err
	}

	go func() {
		for {
			fmt.Println("Waiting here for a new connection")
			conn, err := listener.Accept()
			fmt.Println("Seems like i got a new connection")
			if err != nil {
				fmt.Println("Hello there is an error in accepting a connection")
			}
			var hostname string
			server.connect <- &TcpContext{hostname, conn}
			fmt.Printf("%v", conn)
		}
		done <- "Done"
	}()
	return nil
}

func (server *TcpServer) HandleClientConnect(ctx *TcpContext){
	client := ctx.Conn
	go func(){
		server.handle(ctx)
		server.disconnect <- client
	}()
}

func (server *TcpServer) handle(ctx *TcpContext) (err error){
	clientConn := ctx.Conn
	var backendConn net.Conn
	timeout, _ := time.ParseDuration("2s")
	backendIdleTimeout, _ := time.ParseDuration("10s")
	clientIdleTimeout, _ := time.ParseDuration("10s")
	backendAddress, err := server.scheduler.GetBackend()
	if err != nil {
		fmt.Printf("No backends present to route the request %v\n", err)
		return err
	}
	backendConn, err = net.DialTimeout("tcp", backendAddress, timeout)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Println("Start connection")
	cs := proxy(clientConn, backendConn, backendIdleTimeout)
	bs := proxy(backendConn, clientConn, clientIdleTimeout)

	isTx, isRx := true, true
	//i := 0
	for isTx || isRx {
		select {
		case _, ok := <-cs:
			isRx = ok
		case _, ok2 := <- bs:
			isTx = ok2
		}
		//i += 1
		//fmt.Printf("Hello brother, i am executing in for loop, take me out when connection ends %d\n", i)

	}
	fmt.Println("End connection")
	return nil
}

func (server *TcpServer) HandleClientDisConnect(client net.Conn){
	fmt.Println("Hello not coming here")
	client.Close()
}
