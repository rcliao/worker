package main

import (
	"sync"

	"github.com/rcliao/worker"
)

var wg sync.WaitGroup

func main() {
	workQueue := make(chan worker.Work, 100)
	generatorQuitChan := make(chan bool)
	numWorkers := 5000

	worker.Generator(workQueue, generatorQuitChan)
	worker.Dispatcher(numWorkers, workQueue)

	// to continue processing
	wg.Add(1)
	wg.Wait()
}
