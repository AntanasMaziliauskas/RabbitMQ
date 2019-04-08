package broker

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"time"

	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"github.com/streadway/amqp"
)

type BrokerService interface {
	Init() error
	Stop() error
	GetPersonNode(string, string, types.Person) (types.Person, error)
}

type BrokerClient struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func (b *BrokerClient) Init() error {
	var err error

	rand.Seed(time.Now().UTC().UnixNano())
	b.Conn = &amqp.Connection{}
	b.Ch = &amqp.Channel{}
	b.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Failed to connect to RabbitMQ")

	b.Ch, err = b.Conn.Channel()
	handleError(err, "Failed to open a channel")

	return nil
}

func (b *BrokerClient) Stop() error {
	b.Conn.Close()
	b.Ch.Close()

	return nil
}

func (b *BrokerClient) GetPersonNode(to string, action string, name types.Person) (types.Person, error) {
	res := types.Person{}

	q, err := b.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	handleError(err, "Failed to declare a queue")

	msgs, err := b.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = b.Ch.Publish(
		"",    // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "getPerson"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(name.Name), //Convert
		})
	handleError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			//res.Profession = string(d.Body)
			break
		}
	}

	return res, nil
}

func (b *BrokerClient) GetPersonNode(to string, action string, name types.Person) (types.Person, error) {
	res := types.Person{}

	q, err := b.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	handleError(err, "Failed to declare a queue")

	msgs, err := b.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = b.Ch.Publish(
		"",    // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "getPerson"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(name.Name), //Convert
		})
	handleError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			//res.Profession = string(d.Body)
			break
		}
	}

	return res, nil
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// GobUnmarshal is used internaly in this package to unmarshal gob blob into
// defined interface.
func GobUnmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(v)
}
