package myschedular

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func taskFunction(a string, b int) {
	fmt.Println("THIS IS A JOB", a, b)
}

func GoCronRunner() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			4*time.Second,
		),
		gocron.NewTask(
			taskFunction,
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-time.After(time.Minute):
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		// handle error
	}
}
