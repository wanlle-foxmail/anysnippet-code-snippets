package main

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
)

type JobGuard struct {
	running atomic.Bool
}

// Flow:
//
//	try to flip the running flag from false to true
//	   |
//	   +-> already running -> skip this tick and return ran=false
//	   `-> acquired -> run the job and always clear the flag on exit
func (guard *JobGuard) Run(job func() error) (bool, error) {
	if guard == nil {
		return false, errors.New("guard is required")
	}
	if job == nil {
		return false, errors.New("job is required")
	}
	if !guard.running.CompareAndSwap(false, true) {
		return false, nil
	}
	defer guard.running.Store(false)

	return true, job()
}

func (guard *JobGuard) Running() bool {
	if guard == nil {
		return false
	}
	return guard.running.Load()
}

func main() {
	guard := &JobGuard{}
	ran, err := guard.Run(func() error {
		time.Sleep(20 * time.Millisecond)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ran=%v running=%v", ran, guard.Running())
}
