package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	requests    = 150
	concurrency = 10
	url         = "http://localhost:3000/"
	wait        = sync.WaitGroup{}
)

func main() {
	client := &http.Client{}
	output := make(chan string)

	for i := 0; i < concurrency; i++ {
		go httpWorker(client, output)
	}

	for i := 0; i < requests; i++ {
		output <- url
		fmt.Printf("\r%d hehe ", i+1)
	}
	close(output)

	fmt.Println("")
	os.Stdout.Sync()
	fmt.Println("All URLs handled.")
	os.Stdout.Sync()

	wait.Wait()
	fmt.Println("All workers done.")
}

func httpWorker(client *http.Client, input chan string) {
	// Increment the wait group counter
	wait.Add(1)

	// Iterate over channel input work
	for req := range input {
		if _, err := client.Head(req); err != nil {
			fmt.Println(err)
		}
	}

	// Finish up this worker
	wait.Done()
}
