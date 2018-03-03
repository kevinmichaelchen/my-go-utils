package rabbitmq

import (
	"log"
	"github.com/streadway/amqp"
	"github.com/kevinmichaelchen/my-go-utils"
	"os"
	"fmt"
	"time"
	"encoding/json"
)

// Message represents an empty interface that will be marshaled to JSON and sent as a message.
type Message string //interface{}

// Connection info to our RabbitMQ cluster
var (
	Enabled  = utils.EnvOrBool("RABBITMQ_ENABLED", true)
	User     = os.Getenv("RABBITMQ_USER")
	Password = os.Getenv("RABBITMQ_PASSWORD")
	Host     = os.Getenv("RABBITMQ_HOST")
	Port     = utils.EnvOrInt("RABBITMQ_PORT", 5672)
)

// getConnectionString returns a connection string to our RabbitMQ cluster
func getConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", User, Password, Host, Port)
}

// RabbitListener is an object that will listen for messages.
type RabbitListener struct {
	exchangeName string
	routingKey   string
	queueName    string
	conn         *amqp.Connection
}

// RabbitSender is an object that will send messages.
type RabbitSender struct {
	exchangeName string
	routingKey   string
	conn         *amqp.Connection
}

// NewRabbitSender creates an amqp.Connection, creates the exchange, returns the sender
func NewRabbitSender(exchangeName, routingKey string) *RabbitSender {
	if !Enabled {
		return nil
	}
	conn := createConnection()
	createExchange(conn, exchangeName)
	return &RabbitSender{
		exchangeName: exchangeName,
		routingKey:   routingKey,
		conn:         conn,
	}
}

// NewRabbitListener creates an amqp.Connection, creates the exchange, returns the sender
func NewRabbitListener(exchangeName, routingKey, queueName string) *RabbitListener {
	if !Enabled {
		return nil
	}
	conn := createConnection()
	bindExchangeToQueue(conn, exchangeName, routingKey, queueName)
	return &RabbitListener{
		exchangeName: exchangeName,
		routingKey:   routingKey,
		conn:         conn,
		queueName:    queueName,
	}
}

func (l *RabbitListener) Listen() {
	if !Enabled {
		return
	}

	log.Printf("Opening a channel for %s to listen to %s via key: %s\n", l.queueName, l.exchangeName, l.routingKey)
	ch, err := l.conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Println("Consuming deliveries from channel...")
	msgs, err := ch.Consume(
		l.queueName, // queue
		"",          // consumer
		true,        // auto ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [RECEIVED] %s", d.Body)
		}
	}()

	log.Println("Listening forever...")
	// block forever
	<-forever
}

func (s *RabbitSender) Send(payload interface{}) {
	if !Enabled {
		return
	}

	// TODO reuse channel?
	start := time.Now()
	ch, err := s.conn.Channel()
	log.Printf("Took %s to open a channel...\n", time.Since(start))
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgBytes, err := json.Marshal(payload)
	failOnError(err, "Failed to marshal payload")

	err = ch.Publish(
		s.exchangeName,
		s.routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msgBytes,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s to %s with key: %s", string(msgBytes), s.exchangeName, s.routingKey)
}

// createConnection creates a amqp.Connection, retrying 3 times.
func createConnection() *amqp.Connection {
	var conn *amqp.Connection
	var err error

	connString := getConnectionString()

	log.Printf("Connecting to RabbitMQ URL: %s", connString)

	for i := 0; i < 3; i++ {
		conn, err = amqp.Dial(connString)

		if err != nil {
			log.Printf("Could not connect to RabbitMQ. Will sleep for a bit and then retry")
			time.Sleep(5 * time.Second)
		}
	}

	if conn == nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	} else {
		log.Printf("Successfully connected to RabbitMQ: %s\n", connString)
	}

	return conn
}

// createExchange creates or gets an exchange.
func createExchange(conn *amqp.Connection, exchangeName string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")
}

// bindExchangeToQueue gets/creates the exchange, gets/creates the queue, binds the given exchange to the given queue.
func bindExchangeToQueue(conn *amqp.Connection, exchangeName, routingKey, queueName string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	createExchange(conn, exchangeName)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, exchangeName, routingKey)
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s", msg)
		panic(err)
	}
}
