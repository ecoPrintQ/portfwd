package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	version   string = "1.0.0-alpha"
	buildTime string = "2025-07-21"
	forward   ForwardStruct
	config    ConfigStruct
)

type ForwardStruct struct {
	Protocol string   `yaml:"protocol"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

type ConfigStruct struct {
	Forward             []ForwardStruct `yaml:"forward"`
	ErrorBeforeRecovery int             `yaml:"error_before_recovery" default:"3"`
	TimeBeforeRecovery  int             `yaml:"time_before_recovery" default:"3"`
}

func main() {
	log.Printf("Starting portfwd %s build on %s", version, buildTime)
	configFilePath := os.Getenv("PORTFWD_CONFIG_FILE_PATH")
	if configFilePath == "" {
		configFilePath = "config.yaml"
		log.Printf("PORTFWD_CONFIG_FILE_PATH not defined. Use default configuration file. (config.yaml)")
	}
	log.Printf("Loading configuration file located at %s", configFilePath)
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(config.Forward) == 0 {
		log.Fatal("Nothing to forward! Please check your configuration.")
	}

	log.Printf("-----------------------------------------------")
	log.Printf("Configuration loaded \n %+v", config)
	log.Printf("-----------------------------------------------")

	for index, forward := range config.Forward {
		fmt.Printf("------------Index Forward config %d-------------\n", index)
		log.Printf("Forwarding %+v\n", forward)
		fmt.Printf("-------------------------\n")

		go func(fwd ForwardStruct) {
			var forwardFunc func() error
			var name string

			for indexTo, to := range fwd.To {
				fmt.Printf("------------index To: %d-------------\n", indexTo)
				log.Printf("to %+v\n", to)
				fmt.Printf("-------------------------\n")

				log.Printf("-> Forwarding %s traffic between %s and %s.", fwd.Protocol, fwd.From, to)

				switch fwd.Protocol {
				case "tcp", "tcp4", "tcp6":
					//port := getPort(to)
					//if !isUdpPortInUse(port) {
					forwardFunc = tcpForward(fwd.Protocol, fwd.From, to)
					name = "tcp-forward"
					log.Printf("Forwarding TCP traffic between %s and %s.", fwd.From, to)
					//} else {
					//	fmt.Printf("Error listening on TCP port %d: %v\n", port, err)
					//}
				case "udp", "udp4", "udp6":
					port := getPort(fwd.From)
					if !isUDPPortInUse(fmt.Sprintf("%d", port)) {
						forwardFunc = udpForward(fwd)
						name = "udp-forward"
						//for _, to := range forward.To {
						log.Printf("Forwarding UDP traffic from %s to %s.", fwd.From, to)
						//}
					} else {
						log.Fatalf("Error listening on UDP port %d: %v\n", port, err)
					}
				}

				// Ejecutar el forwarder con el mecanismo de reintento
				runWithRetry(name, forwardFunc)
			}

		}(forward)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}

// runWithRetry ejecuta una función (el forwarder) en un bucle infinito.
// Si retorna un error, espera un momento y vuelve a intentarlo.
func runWithRetry(name string, forwarderFunc func() error) {
	log.Printf("Iniciando servicio: %s\n", name)
	for {
		err := forwarderFunc()
		if err != nil {
			log.Printf("Error en %s: %v. Reintentando en %d segundos...\n", name, err, config.TimeBeforeRecovery)
			time.Sleep(time.Duration(config.TimeBeforeRecovery) * time.Second)
		} else {
			// Esto solo se ejecutaría si la función terminara sin un error.
			// En un forwarder típico, el bucle es interno y no debería llegar aquí.
			//log.Printf("Servicio %s finalizado sin error, terminando bucle de reintento.\n", name)
			break
		}
	}
}

func getPort(addr string) int {
	_, portString, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	port, _ := strconv.Atoi(portString)
	return port
}
