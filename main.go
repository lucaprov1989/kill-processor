package main

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	ps "github.com/mitchellh/go-ps"
)

var programsToKill = []string{"skype", "docker", "compass", "postman", "workbench", "slack", "redis", "mongo", "code helper"}

type program struct {
	pid  int
	name string
}

type result struct {
	name   string
	killed bool
}

type errStruct struct {
	err  error
	name string
}
type programList []program
type resultList []result

func main() {
	start := time.Now()
	ch1 := make(chan *program)
	lookupProcesses := func() {
		ps, _ := ps.Processes()
		// pl := programList{}
		fmt.Printf("Processing %v elements\n", len(ps))
		// 1. for every process get name and pid
		for pp, _ := range ps {
			exec := ps[pp].Executable()
			pid := ps[pp].Pid()
			// 1. for every progam to kill check if its inside the target list
			for i, _ := range programsToKill {
				// TODO: parallelize check
				// TODO: abstract this check in function
				target := programsToKill[i]
				// 1. if process contains programs to kill
				if strings.Contains(strings.ToLower(exec), strings.ToLower(target)) {
					fmt.Printf("Found %v\n", exec)
					pr := program{
						pid,
						exec,
					}
					// 2. send on channel
					ch1 <- &pr
				}
			}
		}
		close(ch1)
	}
	// 1. spin up go routine
	go lookupProcesses()
	errc := make(chan errStruct)
	var res resultList

	// @TODO: get channel inside the function as dependency
	kill := func(pr *program) {
		err := syscall.Kill(pr.pid, 15)

		var killed bool
		if err != nil {
			errToSend := errStruct{err: err, name: pr.name}
			errc <- errToSend
			killed = false
			fmt.Printf("Not killed %v \n", pr.name)
			close(errc)
		} else {
			killed = true
			fmt.Printf("Killed %v \n", pr.name)
		}
		res = append(res, result{name: pr.name, killed: killed})

	}
	// 2. for every process in channel
	for pr := range ch1 {
		// 1. spin up killing go routine
		fmt.Printf("Processing %v\n", pr.name)
		go kill(pr)
	}
	// TODO: abstract this in function
	for e := range errc {
		fmt.Printf("Error during killing of name: %v, err: %v \n", e.name, e.err)
	}
	// processed n. killed n not killed in
	fmt.Printf("Finished Result: %v, in time %v", res, time.Since(start))

}
