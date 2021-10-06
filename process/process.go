package process

import (
	"log"
	"strings"
	"syscall"

	ps "github.com/mitchellh/go-ps"
)

type Program struct {
	pid  int
	name string
}
type Result struct {
	name   string
	killed bool
}

// var programsToKill = []string{"skype", "docker", "compass", "postman", "workbench", "slack", "redis", "mongo", "code helper"}
var programsToKill = []string{"notes", "skype", "app store", "telegram"}

func Lookup(lookUpChannel chan *Program) {
	ps, _ := ps.Processes()
	log.Printf("Processing %v elements\n", len(programsToKill))
	for i := range programsToKill {
		target := programsToKill[i]
		for pp := range ps {
			exec := ps[pp].Executable()
			pid := ps[pp].Pid()
			// 1. if process contains programs to kill
			if strings.Contains(strings.ToLower(exec), strings.ToLower(target)) {
				log.Printf("Found %v\n", exec)
				pr := Program{
					pid,
					exec,
				}
				// 2. send on channel
				lookUpChannel <- &pr
			}
		}
	}
	close(lookUpChannel)
}

func Kill(lookUpChannel chan *Program, killChannel chan *Result) {

	for pr := range lookUpChannel {
		err := syscall.Kill(pr.pid, 15)

		if err != nil {
			log.Printf("Not killed %v \n", pr.name)
			log.Println(err)
			killChannel <- &Result{name: pr.name, killed: false}
		} else {
			log.Printf("Killed %v \n", pr.name)
		}
		killChannel <- &Result{name: pr.name, killed: true}

	}

	close(killChannel)

}
func KillSync(lookUpChannel chan *Program) {

	for pr := range lookUpChannel {
		err := syscall.Kill(pr.pid, 15)

		if err != nil {
			log.Printf("Not killed %v \n", pr.name)
			log.Println(err)
		} else {
			log.Printf("Killed %v \n", pr.name)
		}
	}

}
