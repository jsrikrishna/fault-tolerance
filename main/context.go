package main

import "net"

type Context interface {
	String() string
	Ip() net.IP
	Port() int
}

type TcpContext struct {
	Hostname string
	Conn net.Conn
}

func (t TcpContext) String() string {
	return t.Conn.RemoteAddr().String()
}

func (t TcpContext) Ip() net.IP {
	return t.Conn.RemoteAddr().(*net.TCPAddr).IP
}

func (t TcpContext) Port() int {
	return t.Conn.RemoteAddr().(*net.TCPAddr).Port
}
