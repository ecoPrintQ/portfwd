package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"net"
	"time"
)

func udpForward(forward ForwardStruct) {
	src, err := reuseport.ListenPacket(forward.Protocol, forward.From)
	if err != nil {
		log.Printf("Error. The connection failed (ListenPacket): %v", err)
		return
	}
	defer src.Close()

	dstAddr, err := net.ResolveUDPAddr(forward.Protocol, forward.To[0])
	if err != nil {
		log.Printf("Error resolving destination address: %v\n", err)
		return
	}

	dst, err := net.DialUDP(forward.Protocol, nil, dstAddr)
	if err != nil {
		log.Printf("Error. The connection failed (DialUDP): %v", err)
		return
	}
	defer dst.Close()

	for {
		buf := make([]byte, 2048)
		n, addr, err := src.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading from UDP socket: %v\n", err)
		}
		//log.Printf(fmt.Sprintf("%v", buf))
		log.Printf(fmt.Sprintf("Package received from %v: %s\n", addr, string(buf[:n])))

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Printf("Failed to resend packet: %v", err)
			continue
		}

		fmt.Printf("Package forwarded to  %v\n", dstAddr)

		// Esperar la respuesta del remoteAddr
		// Leemos la respuesta del remoteConn (la cual puede ser un paquete UDP enviado de vuelta)
		// NOTA: Esto podría ser bloqueante si no hay respuesta
		dst.SetReadDeadline(time.Now().Add(5 * time.Second)) // Timeout de 5 segundos
		n, err = dst.Read(buf)
		if err != nil {
			log.Printf("Error receiving response from remoteAddr: %v", err)
			continue
		}

		// Mostrar la respuesta recibida del remoteAddr
		fmt.Printf("Response received from %v: %s\n", dstAddr, string(buf[:n]))

		// Reenviar la respuesta al cliente original
		_, err = src.WriteTo(buf[:n], addr)
		if err != nil {
			log.Printf("Error when resending the response to the customer: %v", err)
			continue
		}

		fmt.Printf("Response forwarded to customer %v\n", addr)

	}
}
