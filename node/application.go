package node

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/AntanasMaziliauskas/RabbitMQ/node/person"
	"github.com/streadway/amqp"
)

type Application struct {
	Tasks  map[string]interface{}
	Name   string
	Person person.PersonService
}

func (a *Application) Init() {
	//Tasku sarasas
	//app := Server{Name: *name,
	a.Tasks = map[string]interface{}{
		"getPerson": a.GetPerson,
	}
	//duombaze uzkuriam
	a.Person.Init()
	//Klausomes kanalo
	a.Listen()

}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (a *Application) Listen() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		a.Name, // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	handleError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	handleError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//log.Printf("Received a message Header: %s", d.Headers["action"])
			a.Tasks[d.Headers["action"].(string)].(func(amqp.Delivery, amqp.Channel))(d, *ch)
		}
	}()

	log.Printf("Waiting for messages sent to %s...", a.Name)
	<-forever
}

func (a *Application) GetPerson(d amqp.Delivery, ch amqp.Channel) {
	response := a.Person.GetPerson(string(d.Body))
	//response := types.Person{Name: "Jonas", Profession: "bilekas", Age: 22}
	atsakas, err := GobMarshal(response)
	handleError(err, "Failed to Marshal with gob")

	err = ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          atsakas,
		})
	d.Ack(false)
	handleError(err, "Failed to publish a message")
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

// GobUnmarshal is used internaly in this package to unmarshal gob blob into
// defined interface.
func GobUnmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return gob.NewDecoder(b).Decode(v)
}
