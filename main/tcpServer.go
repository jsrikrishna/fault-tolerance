package main

import (
	"net"
	"fmt"
	"fault-tolerance/config"
	"fault-tolerance/ping"
)

type TcpServer struct {
	name                     string
	bindTo                   string
	clientIdleTimeout        string
	backendIdleTimeout       string
	backendConnectionTimeout string
	listener                 net.TCPListener
	connect                  chan (*TcpContext)
	disconnect               chan (net.Conn)
	done                     chan string
	scheduler                *ping.Scheduler
}

func New(configuration config.Configuration, scheduler *ping.Scheduler, done chan string) *TcpServer {
	tcpServer := &TcpServer{
		name: configuration.Name,
		bindTo:configuration.BindTo,
		clientIdleTimeout: configuration.ClientIdleTimeout,
		backendIdleTimeout: configuration.BackendIdleTimeout,
		backendConnectionTimeout: configuration.BackendConnectionTimeout,
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
			fmt.Printf("Connection is %v\n", conn)
		}
		done <- "Done"
	}()
	return nil
}

func (server *TcpServer) HandleClientConnect(ctx *TcpContext) {
	client := ctx.Conn
	go func() {
		server.handle(ctx)
		server.disconnect <- client
	}()
}

func (server *TcpServer) handle(ctx *TcpContext) (err error) {
	clientConn := ctx.Conn
	var backendConn net.Conn
	backendAddress, err := server.scheduler.GetBackend()
	if err != nil {
		fmt.Printf("No backends present to route the request %v\n", err)
		return err
	}
	backendConn, err = net.DialTimeout("tcp", backendAddress, ParseDurationOrDefault(server.backendConnectionTimeout, 0))
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Println("Start connection")
	cs := proxy(clientConn, backendConn, ParseDurationOrDefault(server.backendIdleTimeout, 0))
	bs := proxy(backendConn, clientConn, ParseDurationOrDefault(server.clientIdleTimeout, 0))

	isTx, isRx := true, true
	for isTx || isRx {
		select {
		case _, ok := <-cs:
			isRx = ok
		case _, ok2 := <-bs:
			isTx = ok2
		}
	}
	fmt.Println("End connection")
	return nil
}

func (server *TcpServer) HandleClientDisConnect(client net.Conn) {
	fmt.Println("Closing Client Connection now")
	client.Close()
}
