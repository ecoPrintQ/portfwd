package main

import (
	"github.com/libp2p/go-reuseport"
	"io"
	"log"
	"net"
)

func tcpForward(protocol string, from string, to string) {
	listener, err := reuseport.Listen(protocol, to)

	if err != nil {
		log.Printf("The connection failed: %v", err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("The connection was not accepted: %v", err)
		}

		client, err := net.Dial(protocol, from)

		if err != nil {
			log.Printf("The connection failed: %v", err)
			conn.Close()
			continue
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
