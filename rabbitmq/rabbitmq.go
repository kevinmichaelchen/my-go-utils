package rabbitmq

import "log"

func SendMessage(s Message) {
	go func(ss Message) {
		MessageQueue <- ss
	}(s)
}

// Message represents an empty interface that will be marshaled to JSON and sent as a message.
type Message string //interface{}

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

// MessageQueue is a buffered channel that we can send messages on.
// TODO customize channel capacity
var MessageQueue = make(chan Message, 10)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Message
}

func NewDispatcher(numWorkers int) *Dispatcher {
	pool := make(chan chan Message, numWorkers)
	d := &Dispatcher{WorkerPool: pool}
	d.run(numWorkers)
	return d
}

func (d *Dispatcher) run(numWorkers int) {
	for i := 1; i <= numWorkers; i++ {
		worker := NewWorker(i, d.WorkerPool)
		worker.Start()
	}

	if numWorkers > 0 {
		go d.dispatch()
	}
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case message := <-MessageQueue:
			// a message request has been received
			go func(message Message) {
				// try to obtain a worker message channel that is available.
				// this will block until a worker is idle
				messageChannel := <-d.WorkerPool

				// dispatch the message to the worker message channel
				messageChannel <- message
			}(message)
		}
	}
}
