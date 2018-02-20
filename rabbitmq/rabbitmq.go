package rabbitmq

func SendMessage(s Message) {
	go func(ss Message) {
		MessageQueue <- ss
	}(s)
}

// Message represents an empty interface that will be marshaled to JSON and sent as a message.
type Message string //interface{}
