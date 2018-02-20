package rabbitmq

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Message
}

// MessageQueue is a buffered channel that we can send messages on.
var MessageQueue chan Message

func NewDispatcher(numWorkers, channelCapacity int) *Dispatcher {
	if channelCapacity > 0 {
		MessageQueue = make(chan Message, channelCapacity)
	} else {
		MessageQueue = make(chan Message, 10)
	}
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
