package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
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
	Forward []ForwardStruct `yaml:"forward"`
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

	for _, forward = range config.Forward {
		switch forward.Protocol {
		case "tcp", "tcp4", "tcp6":
			for _, to := range forward.To {
				go tcpForward(forward.Protocol, forward.From, to)
				log.Printf("Forwarding TCP traffic between %s and %s.", forward.From, to)
			}
		case "udp", "udp4", "udp6":
			go udpForward(forward)
			for _, to := range forward.To {
				log.Printf("Forwarding UDP traffic from %s to %s.", forward.From, to)
			}
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
