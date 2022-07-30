package ampq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type queue struct {
	url  string
	name string

	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	closed       bool

	consumers []messageConsumer
}

type messageConsumer func(string)

func NewQueue(url string, qName string) *queue {
	q := new(queue)
	q.url = url
	q.name = qName

	q.connect()
	go q.reconnector()

	return q
}

func (q *queue) Send(message []byte) {
	err := q.channel.Publish(
		"",     // exchange
		q.name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		log.Println("Sending message to queue failed: ", err)
	}
}

func (q *queue) Consumer() (<-chan amqp.Delivery, error) {
	log.Println("Registering consumer...")
	deliveries, err := q.registerQueueConsumer()
	if err != nil {
		log.Println("Consumer registration failed: ", err)
	}
	log.Println("Consumer registered!")
	return deliveries, nil
}

func (q *queue) connect() {
	for {
		conn, err := amqp.Dial(q.url)
		if err == nil {
			q.connection = conn
			q.errorChannel = make(chan *amqp.Error)
			q.connection.NotifyClose(q.errorChannel)

			q.openChannel()
			q.declareQueue()

			return
		}

		if err != nil {
			log.Println("Connection to rabbitmq failed. Retrying in 1 sec... ", err)
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func (q *queue) reconnector() {
	for {
		err := <-q.errorChannel
		if !q.closed {
			log.Println("Reconnecting after connection closed: ", err)

			q.connect()
		}
	}
}

func (q *queue) declareQueue() {
	_, err := q.channel.QueueDeclare(
		q.name, // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Println("Queue declaration failed: ", err)
	}
}

func (q *queue) openChannel() {
	channel, err := q.connection.Channel()
	if err != nil {
		log.Println("Opening channel failed: ", err)
	}
	q.channel = channel
}

func (q *queue) registerQueueConsumer() (<-chan amqp.Delivery, error) {
	msgs, err := q.channel.Consume(
		q.name, // queue
		"",     // messageConsumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println("Consuming messages from queue failed: ", err)
	}
	return msgs, err
}
