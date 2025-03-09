package main

import (
	"net"
)

// customListener wraps a TCPListener
type customListener struct {
	*net.TCPListener
}

// Accept waits for and returns the next connection to the listener
func (ln *customListener) Accept() (net.Conn, error) {
	conn, err := ln.TCPListener.Accept()
	if err != nil {
		return conn, err
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if ok {
		tcpConn.SetNoDelay(true)
	}

	return conn, err
}
