package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()
	totalJobs := 5
	jobs := make(chan int, totalJobs)
	var wg sync.WaitGroup

	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	for job := 1; job <= totalJobs; job++ {
		jobs <- job
	}

	close(jobs)
	wg.Wait()
	fmt.Println("Total time ", time.Since(startTime))
}

func worker(w int, jobs chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		processJobs(w, job)
	}
}

func processJobs(w int, job int) {
	fmt.Println("Worker", w, "started  job", job)
	time.Sleep(time.Second)
	fmt.Println("Worker", w, "finished job", job)
}
