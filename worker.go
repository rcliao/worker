package worker

import (
	"fmt"
	"math/rand"
	"time"
)

type Work struct {
	Data int32
}

type Worker struct {
	ID          int
	Work        chan Work
	WorkerQueue chan *Worker
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan *Worker) Worker {
	return Worker{
		ID:          id,
		Work:        make(chan Work),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
}

func (w *Worker) Run() {
	go func() {
		for {
			w.WorkerQueue <- w

			select {
			case work := <-w.Work:
				fmt.Println("Worker", w.ID, "received work", work)
				time.Sleep(time.Second)
				result := work.Data * 2
				fmt.Println("Worker", w.ID, "finished work", result)
			case <-w.QuitChan:
				fmt.Println("Worker", w.ID, "stopped")
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func Generator(workQueue chan Work, quit chan bool) {
	ticket := time.NewTicker(1 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticket.C:
				// send work
				i := rand.Int31n(100)
				workQueue <- Work{i}
			case <-quit:
				ticket.Stop()
				return
			}
		}
	}()
}

func Dispatcher(numWorkers int, workQueue chan Work) {
	workerQueue := make(chan *Worker, numWorkers)

	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, workerQueue)
		fmt.Println("Starting worker", i)
		worker.Run()
	}

	go func() {
		for {
			select {
			case work := <-workQueue:
				fmt.Println("Dispatcher receives work", work)
				go func() {
					worker := <-workerQueue
					fmt.Println("Dispatching work")
					worker.Work <- work
				}()
			}
		}
	}()
}
