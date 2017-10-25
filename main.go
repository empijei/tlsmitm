package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var rules []listener

var nolog = flag.Bool("silent", false, "do not log the traffic")

//TODO print example
var conf = flag.String("conf", "conf.json", "a conf file")

func main() {
	flag.Parse()
	//TODO parse CLI parameters and create rules
	f, err := os.Open(*conf)
	if err != nil {
		log.Printf("Conf file not found: %s\n", *conf)
		return
	}
	log.Println("Loading conf...")
	err = loadConf(f)
	if err != nil {
		log.Println("Error while loading conf: ", err)
		return
	}
	log.Println("Conf loaded...")
	for _, l := range rules {
		l := l
		go l.Listen()
	}

	//Wait for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
