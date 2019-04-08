package main

import (
	"flag"
	"log"

	"github.com/streadway/amqp"
)

/*
type Person struct {
	Name       string `bson:"name"`
	Age        int64  `bson:"age"`
	Profession string `bson:"profession"`
}
*/
func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	to := flag.String("topic", "Node01", "Topic to be used. # by default.")
	//action = flag.String("action", "", "")
	msg := flag.String("msg", "Hi", "Message to be sent.")
	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"topics", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	handleError(err, "Failed to declare an exchange")

	err = ch.Publish(
		"topics", // exchange
		*to,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(*msg),
		})
	handleError(err, "Failed to publish a message")

	log.Printf("Sending message to %s: %s", *to, *msg)

	//JUDAM TOLIAU
	result, err := broker.Call("node-id", "getPerson", message.GetPerson{Name: "jonas"})

}

//GOB example
/*
import (
	"bytes"
	"encoding/gob"
)

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


*/
