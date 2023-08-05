// Copied from a reddit post written by /u/jerf (Thanks Jerf)
// https://www.reddit.com/r/golang/comments/947jul/best_practises_for_pool_of_go_routines/

package main

import (
	"fmt"
	"sync"
	"time"
)

/*
In many languages, you need a "pool" mechanism to handle threads, because
threads are expensive and you don't want to start many of them, so once
started, you want to get as much value out of them as possible. Since
that can be a pain, you often use a library to make it easier. People
coming to go from other languages then wonder where the pool libraries
are in Go.
The answer is that while Go, strictly speaking, does not have a built-in
"pool", the primitives that it does support are so close to what you need
that there is not much room for a library to help you out. This code snippet
will demonstrate how to create a "worker pool" of goroutines, dispatch jobs
to those pools, shut them down properly, and then continue on. At the end,
I'll comment on the things to watch out for when using this technique.
*/

const (
	NumberOfWorkers = 3
)

func main() {
	// make the tasks channel
	tasks := make(chan string)

	wg := sync.WaitGroup{}
	// you often see people adding them one by one as they spawn goroutines
	// but you can add them in one shot if you know in advance.
	wg.Add(NumberOfWorkers)

	for i := 0; i < NumberOfWorkers; i++ {
		go func(workerNum int) {
			// in real code, a defer function to recover is
			// a good idea here, because any panics would
			// otherwise crash the entire program
			defer wg.Done()

			for {
				// pull tasks from queue until done
				task, ok := <-tasks
				if !ok {
					// we're done
					return
				}
				fmt.Println("Worker", workerNum, ":", task)
				// JUST FOR THIS DEMO, give the other workers a
				// chance to catch jobs; normally this is unnecessary
				// of course!
				time.Sleep(time.Millisecond)
			}
		}(i) // in Go, if you try to close on the i above, you'll have a race
	}

	// we now have NumberOfWorkers running. Give them tasks:
	for i := 0; i < 10; i++ {
		tasks <- fmt.Sprintf("This is task %d", i)
	}

	// Signal that we're done with our work:
	close(tasks)

	// Wait for the tasks to complete:
	wg.Wait()

	// You're done!
}

/*
It is tempting to say "But I could use a library for some of those things
up there". However, look at what you're actually saving. Is it really an
advantage to type "Pool.EndJobs()" vs. "close(taskChan)"? Every Go
programmer will understand the latter, and its precise semantics. The
former? Not so much... does it immediately terminate the pool or wait for
jobs to finish. Does the call synchronously wait for the jobs to finish?
Is there any sort of context involvement? You'll have to document this
all in your new library, and people will have to learn it, and it won't
carry over to somebody else's library, whereas `close(taskChan)` is
completely obvious.
The biggest traps I see in this pattern are:
1. Watch out for the amount of work your workers are doing. The pool may
   be cheap, but it is still not free. You want your workers to be doing
   enough work that the coordination of spinning up goroutines and using
   channels is a negligible fraction of the time. Something like "printing
   a single string" or "adding two numbers together" (a common beginner
   test task) is too small. (A simple solution is to be sure to bundle up
   enough work in one message to make it worth it.)
   However, generally, if your work tasks are so small that it's too
   expensive to spin up some worker goroutines, they're too expensive to
   be dispatching across cores via *any* mechanism. Don't underestimate
   coordination costs; for small tasks it can be fastest to just do them
   on one core in one goroutine regardless.
2. Certain request patterns could make this problematic; if you have
   *highly* bursty requests, then you may grind your process to a halt
   trying to spin up a lot of goroutines for each request when a pool
   could have worked out. In that case, you still don't need a library
   per se; you just take the above code and spin the pool up once for
   that task. There still isn't much use for a "generic worker" pool in
   Go; it's cleaner code and except in rare cases, not that much more
   resource-expensive to just go ahead and spin up a pool per task,
   which keeps the tasks from being coupled to each other.
*/
