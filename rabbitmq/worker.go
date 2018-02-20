package rabbitmq

import "log"

// Worker represents the worker that sends the message.
type Worker struct {
	WorkerPool     chan chan Message
	MessageChannel chan Message
	quit           chan bool
	id             int
}

func NewWorker(id int, workerPool chan chan Message) Worker {
	return Worker{
		WorkerPool:     workerPool,
		MessageChannel: make(chan Message),
		quit:           make(chan bool),
		id:             id}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			log.Printf("Registering worker %d into worker queue...\n", w.id)
			w.WorkerPool <- w.MessageChannel

			select {
			case message := <-w.MessageChannel:
				// we have received a work request.
				log.Printf("Worker %d received message: %s", w.id, message)

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
