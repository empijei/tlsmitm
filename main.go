package main

import (
	"os"
	"os/signal"
)

var rules []listener

func main() {

	//TODO parse parameters and create rules

	rules = append(rules, listener{
		localport:  ":9443",
		remoteport: ":8000",
		remoteip:   "127.0.0.1",
	})

	for _, l := range rules {
		l := l
		go l.Listen()
	}

	//Wait for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
