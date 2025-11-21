package main

import (
	"fmt"
	"github.com/libp2p/go-reuseport"
	"log"
	"net"
	"time"
)

func udpForward(forward ForwardStruct) func() error {
	return func() error {
		log.Printf("New forward UDP: %s", forward.From)
		var errorCount = 0
		src, err := reuseport.ListenPacket(forward.Protocol, forward.From)
		if err != nil {
			errF := fmt.Errorf("Error. The connection failed (ListenPacket): %v", err)
			log.Fatalf("%v", errF)
		}
		defer src.Close()

		dstAddr, err := net.ResolveUDPAddr(forward.Protocol, forward.To[0])
		if err != nil {
			errF := fmt.Errorf("Error resolving destination address: %v\n", err)
			log.Fatalf("%v", errF)
		}

		dst, err := net.DialUDP(forward.Protocol, nil, dstAddr)
		if err != nil {

			errF := fmt.Errorf("Error. The connection failed (DialUDP): %v", err)
			log.Fatalf("%v", errF)
		}
		defer dst.Close()

		for {
			if errorCount >= config.ErrorBeforeRecovery {
				errF := fmt.Errorf("Too many errors (%d). Restart", errorCount)
				return errF
			}

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
				errorCount += 1
				continue
			}

			log.Printf("Package forwarded to  %v\n", dstAddr)

			// Esperar la respuesta del remoteAddr
			// Leemos la respuesta del remoteConn (la cual puede ser un paquete UDP enviado de vuelta)
			// NOTA: Esto podría ser bloqueante si no hay respuesta
			dst.SetReadDeadline(time.Now().Add(5 * time.Second)) // Timeout de 5 segundos
			n, err = dst.Read(buf)
			if err != nil {
				log.Printf("Error receiving response from remoteAddr: %v", err)
				errorCount += 1
				continue
			}

			// Mostrar la respuesta recibida del remoteAddr
			log.Printf("Response received from %v: %s\n", dstAddr, string(buf[:n]))

			// Reenviar la respuesta al cliente original
			_, err = src.WriteTo(buf[:n], addr)
			if err != nil {
				log.Printf("Error when resending the response to the customer: %v", err)
				errorCount += 1
				continue
			}

			log.Printf("Response forwarded to customer %v\n", addr)
		}
	}
}

func isUDPPortInUse(port string) bool {
	log.Printf("Checking if port %v is in use\n", port)
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Printf("Port %v not in use\n", port)
		return true
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("Port %v is in use\n", port)
		// El puerto está en uso o no se puede acceder
		return true
	}
	conn.Close()
	log.Printf("Port %v is not in use\n", port)
	return false
}

func isUdpPortInUsec(port int) bool {
	log.Printf("------------------------Check UDP port %d\n", port)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Error resolving UDP address for port %d: %v\n", port, err)
		return true // Treat as in use if address resolution fails
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {

		// If an error occurs, it likely means the port is in use
		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" && opErr.Err.Error() == "address already in use" {
			return true
		}
		fmt.Printf("Error listening on UDP port %d: %v\n", port, err)
		return true // Treat other errors as in use as well
	}
	defer conn.Close() // Close the listener if successful
	fmt.Printf("Listen sin errores- ------------------------------------------->")
	return false // Port is not in use

}
