package worker

import (
	"fmt"
	"time"
)

// Work is a function to process for job with certain payload
type Process func(payload string) string

// Job represents a unit of work with Fn to do the work with Payload
type Job struct {
	Payload string
	Process Process
}

// Worker is a unit of individual worker to do job
type Worker struct {
	ID          int
	Jobs        chan Job
	WorkerQueue chan *Worker
	QuitChan    chan bool
}

// NewWorker is a convient constructor for creating a new Worker
func NewWorker(id int, workerQueue chan *Worker) Worker {
	return Worker{
		ID:          id,
		Jobs:        make(chan Job),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
}

// Run is a main method to start a bunch of worker and start doing jobs
func (w *Worker) Run() {
	go func() {
		for {
			w.WorkerQueue <- w

			select {
			case job := <-w.Jobs:
				fmt.Println("Worker", w.ID, "received job", job)
				time.Sleep(time.Second)
				result := job.Process(job.Payload)
				fmt.Println("Worker", w.ID, "finished job", result)
			case <-w.QuitChan:
				fmt.Println("Worker", w.ID, "stopped")
				return
			}
		}
	}()
}

// Stop stops the Run method that start a bunch of jobs
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// Dispatcher dispatches works to individual worker
func Dispatcher(numWorkers int, workQueue chan Job) {
	workerQueue := make(chan *Worker, numWorkers)

	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, workerQueue)
		fmt.Println("Starting worker", i)
		worker.Run()
	}

	go func() {
		for {
			select {
			case job := <-workQueue:
				fmt.Println("Dispatcher receives work", job)
				go func() {
					worker := <-workerQueue
					fmt.Println("Dispatching work")
					worker.Jobs <- job
				}()
			}
		}
	}()
}
