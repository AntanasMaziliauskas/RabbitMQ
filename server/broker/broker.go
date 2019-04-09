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
	UpsertPerson(string, string, types.Person) error
	ListPersonsBroadcast() error
}

type BrokerClient struct {
	Conn    *amqp.Connection
	Ch      *amqp.Channel
	Consume <-chan amqp.Delivery
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

	err = b.Ch.ExchangeDeclare(
		"nodesdirects", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	handleError(err, "Failed to declare an exchange")

	q, err := b.Ch.QueueDeclare(
		"response", // name
		false,      // durable
		false,      // delete when usused
		false,      // exclusive
		false,      // noWait
		nil,        // arguments
	)
	handleError(err, "Failed to declare a queue")

	b.Consume, err = b.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")
	return nil
}

func (b *BrokerClient) Stop() error {
	b.Conn.Close()
	b.Ch.Close()

	return nil
}

func (b *BrokerClient) GetPersonNode(to string, action string, person types.Person) (types.Person, error) {
	res := types.Person{}

	corrId := randomString(32)
	per, _ := GobMarshal(person)

	log.Println(to)
	//to = "nodes." + to
	err := b.Ch.Publish(
		"nodesdirects", // exchange
		to,             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "getPerson"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       "response",  //q.Name
			Body:          []byte(per), //Convert
		})
	handleError(err, "Failed to publish a message")

	for d := range b.Consume {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			//res.Profession = string(d.Body)
			break
		}
	}

	return res, nil
}

func (b *BrokerClient) UpsertPerson(to string, action string, person types.Person) error {
	res := types.Person{}

	corrId := randomString(32)
	per, _ := GobMarshal(person)

	err := b.Ch.Publish(
		"",    // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "upsertPerson"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       "response",
			Body:          []byte(per), //Convert
		})
	handleError(err, "Failed to publish a message")

	for d := range b.Consume {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			//res.Profession = string(d.Body)
			break
		}
	}

	return nil
}

//TODO: FANOUT
func (b *BrokerClient) ListPersonsBroadcast() error {
	res := []types.Person{}
	result := []types.Person{}

	/*err := b.Ch.ExchangeDeclare(
		"broadcast", // name
		"fanout",    // type
		false,       // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	handleError(err, "Failed to declare exchange")*/

	/*err = b.Ch.QueueBind(
		q.Name,    // queue name
		"nodes.*", // routing key
		"tasks",   // exchange
		false,
		nil,
	)
	handleError(err, "Failed to declare binding")*/

	corrId := randomString(32)
	//per, _ := GobMarshal(person)

	err := b.Ch.Publish(
		"nodes", // exchange
		"nodes", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "listPersons"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       "response",
			//Body:          []byte(per), //Convert
		})
	handleError(err, "Failed to publish a message")

	//Kaip cia viskas atrodys?
	for d := range b.Consume {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			result = append(result, res...)
			//res.Profession = string(d.Body)
			break
		}
	}

	return nil
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

// GobMarshal is used to marshal struct into gob blob
func GobMarshal(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
