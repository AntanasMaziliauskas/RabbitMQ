package node

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/AntanasMaziliauskas/RabbitMQ/node/person"
	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"github.com/streadway/amqp"
)

type Application struct {
	Tasks   map[string]interface{}
	Name    string
	Person  person.PersonService
	Conn    *amqp.Connection
	Ch      *amqp.Channel
	Command <-chan amqp.Delivery
}

func (a *Application) Init() {
	//Tasku sarasas
	//app := Server{Name: *name,
	a.Tasks = map[string]interface{}{
		"getPerson":    a.GetPerson,
		"upsertPerson": a.UpsertPerson,
		"listPersons":  a.ListPersons,
	}
	//duombaze uzkuriam
	a.Person.Init()
	//Setinam eiles ir kanala
	a.Rabbit()
	//Klausomes kanalo
	a.Listen()

}

/*
func (a *Application) Greet() {


}*/

func (a *Application) Rabbit() {
	var err error

	a.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	a.Ch, err = a.Conn.Channel()
	handleError(err, "Failed to open a channel")
	//defer ch.Close()

	q, err := a.Ch.QueueDeclare(
		"comman", // name
		true,     // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	handleError(err, "Failed to declare a queue")

	err = a.Ch.ExchangeDeclare(
		"nodesdirects", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	handleError(err, "Failed to declare an exchange")

	log.Println(a.Name)
	//routing := "nodes." + a.Name
	err = a.Ch.QueueBind(
		q.Name,         // queue name
		a.Name,         // routing key
		"nodesdirects", // exchange
		false,
		nil)
	handleError(err, "Failed to bind a queue")

	a.Command, err = a.Ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")
}

func (a *Application) Stop() {
	a.Conn.Close()
	a.Ch.Close()
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (a *Application) Listen() {

	forever := make(chan bool)

	go func() {
		for d := range a.Command {
			//if d.Headers["action"] != nil {
			log.Printf("Received a message: %s", d.Headers["action"])
			//log.Printf("Received a message Header: %s", d.Headers["action"])
			a.Tasks[d.Headers["action"].(string)].(func(amqp.Delivery))(d)
			//}
		}
	}()

	log.Printf("Waiting for messages sent to %s...", a.Name)
	<-forever
}

func (a *Application) GetPerson(d amqp.Delivery) {
	var person *types.Person
	GobUnmarshal(d.Body, &person)
	//str := person.ID.Hex()
	response := a.Person.GetPerson(person.ID)
	//response := types.Person{Name: "Jonas", Profession: "bilekas", Age: 22}
	atsakas, err := GobMarshal(response)
	handleError(err, "Failed to Marshal with gob")

	err = a.Ch.Publish(
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

func (a *Application) UpsertPerson(d amqp.Delivery) {
	var person *types.Person

	GobUnmarshal(d.Body, &person)
	a.Person.UpsertPerson(person)
	//response := types.Person{Name: "Jonas", Profession: "bilekas", Age: 22}
	//atsakas, err := GobMarshal(response)
	//handleError(err, "Failed to Marshal with gob")

	err := a.Ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          []byte("Done"),
		})
	d.Ack(false)
	handleError(err, "Failed to publish a message")
}

func (a *Application) ListPersons(d amqp.Delivery) {
	//var person *types.Person
	//GobUnmarshal(d.Body, &person)
	//str := person.ID.Hex()
	response, _ := a.Person.ListPersons()
	//response := types.Person{Name: "Jonas", Profession: "bilekas", Age: 22}
	atsakas, err := GobMarshal(response)
	handleError(err, "Failed to Marshal with gob")

	err = a.Ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          atsakas,
		})
	if err == nil {
		log.Println("Reply sent")
	}
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
