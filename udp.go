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
		log.Printf("The connection failed (ListenPacket): %v", err)
	}
	defer src.Close()

	dstAddr, err := net.ResolveUDPAddr(forward.Protocol, forward.To[0])
	if err != nil {
		log.Printf("Error resolving destination address: %v\n", err)
	}

	dst, err := net.DialUDP(forward.Protocol, nil, dstAddr)
	if err != nil {
		log.Printf("The connection failed (DialUDP): %v", err)
	}
	defer dst.Close()

	for {
		log.Printf("For...")
		buf := make([]byte, 2048)
		n, addr, err := src.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading from UDP socket: %v\n", err)
		}
		//log.Printf(fmt.Sprintf("%v", buf))
		log.Printf(fmt.Sprintf("Paquete recibido de %v: %s\n", addr, string(buf[:n])))

		_, err = dst.Write(buf[:n])
		if err != nil {
			log.Printf("Error al reenviar el paquete: %v", err)
			continue
		}

		fmt.Printf("Paquete reenviado a %v\n", dstAddr)

		// Esperar la respuesta del remoteAddr
		// Leemos la respuesta del remoteConn (la cual puede ser un paquete UDP enviado de vuelta)
		// NOTA: Esto podría ser bloqueante si no hay respuesta
		dst.SetReadDeadline(time.Now().Add(5 * time.Second)) // Timeout de 5 segundos
		n, err = dst.Read(buf)
		if err != nil {
			log.Printf("Error al recibir respuesta del remoteAddr: %v", err)
			continue
		}

		// Mostrar la respuesta recibida del remoteAddr
		fmt.Printf("Respuesta recibida de %v: %s\n", dstAddr, string(buf[:n]))

		// Reenviar la respuesta al cliente original
		_, err = src.WriteTo(buf[:n], addr)
		if err != nil {
			log.Printf("Error al reenviar la respuesta al cliente: %v", err)
			continue
		}

		fmt.Printf("Respuesta reenviada al cliente %v\n", addr)

	}
}
