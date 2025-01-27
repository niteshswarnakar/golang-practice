package workerpool

import (
	"fmt"
	"sync"
	"time"
)

func Worker(job string, worker int) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Performing Job : ", job, " by Worker : ", worker)
}

func WorkerPool() {
	jobs := []string{"job1", "job2", "job3", "job4", "job5", "job6", "job7", "job8", "job9", "job10"}
	workers := 4

	wg := &sync.WaitGroup{}

	jobChannel := make(chan string, len(jobs))

	for _, job := range jobs {
		jobChannel <- job
	}
	close(jobChannel)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		i := i
		go func(worker int) {
			defer wg.Done()
			for job := range jobChannel {
				Worker(job, i+1)
			}
		}(i)
	}

	wg.Wait()

}
