package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/rcliao/worker"
)

var wg sync.WaitGroup

func main() {
	workQueue := make(chan worker.Job, 100)
	generatorQuitChan := make(chan bool)
	numWorkers := 5000

	generator(workQueue, generatorQuitChan)
	worker.Dispatcher(numWorkers, workQueue)

	// to continue processing
	wg.Add(1)
	wg.Wait()
}

func generator(workQueue chan worker.Job, quit chan bool) {
	ticket := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticket.C:
				// send work
				i := rand.Int31n(100)
				workQueue <- worker.Job{
					fmt.Sprint(i),
					func(payload string) string {
						result, _ := strconv.Atoi(payload)
						return fmt.Sprint(result * 2)
					},
				}
			case <-quit:
				ticket.Stop()
				return
			}
		}
	}()
}
