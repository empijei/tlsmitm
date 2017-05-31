package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
)

var rules []listener

func main() {

	//TODO parse CLI parameters and create rules

	//TODO allow TLS -> Plain and Plain -> TLS

	//This is the creation of a sample rule
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	rules = append(rules, listener{
		localport:   ":9443",
		remoteport:  ":8000",
		remoteip:    "127.0.0.1",
		secure:      true,
		protoSwitch: false,
		certconf: &tls.Config{
			Certificates:       []tls.Certificate{cer},
			InsecureSkipVerify: true,
		},
	})
	//Rule creation ends here

	for _, l := range rules {
		l := l
		go l.Listen()
	}

	//Wait for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
