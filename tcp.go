package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"io"
	"net"
)

func tcpForward(protocol string, from string, to string) func() error {
	return func() error {
		listener, err := reuseport.Listen(protocol, to)
		if err != nil {
			errF := fmt.Errorf("The connection failed: %v", err)
			return errF
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				errF := fmt.Errorf("The connection was not accepted: %v", err)
				return errF
			}

			client, err := net.Dial(protocol, from)
			if err != nil {
				conn.Close()
				errF := fmt.Errorf("The connection failed: %v", err)
				return errF
			}

			go func() {
				defer client.Close()
				defer conn.Close()
				io.Copy(client, conn)
			}()

			go func() {
				defer client.Close()
				defer conn.Close()
				io.Copy(conn, client)
			}()
		}
	}
}
