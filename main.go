package main

import (
	"log"
	p "lookupProcesses/process"
	"time"
)

func main() {
	start := time.Now()
	lookUpChannel := make(chan *p.Program)
	killChannel := make(chan *p.Result)
	// 1. spin up go routines to:
	// lookup processes to kill and send to the channel
	go p.Lookup(lookUpChannel)
	// read from the channel, kill and send the result to channel
	go p.Kill(lookUpChannel, killChannel)
	for res := range killChannel {
		log.Printf("Program killed: %v", res)
	}
	log.Printf("Finished in time %v", time.Since(start))
}
