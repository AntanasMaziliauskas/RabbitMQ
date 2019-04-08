package node

/*
import (
	"flag"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/streadway/amqp"
)

type Person struct {
	Name       string `bson:"name"`
	Age        int64  `bson:"age"`
	Profession string `bson:"profession"`
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	name := flag.String("name", "Node01", "Name of the Node")
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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	handleError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,   // queue name
		*name,    // routing key
		"topics", // exchange
		false,
		nil)
	handleError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
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
		}
	}()

	log.Printf("Waiting for messages sent to %s...", *name)
	<-forever
}

func GetPerson(p Packet, q *Queue) {
	payload := message.GetPerson{}
	gob.Unmarshal(p.Body, &payload)

	person := fromDatabase.Find(bson.M{"name", payload.Name})
	answer := gob.Marshal(messages.ResponsePerson{Name: person.Name})

	q.Publish(answer)
}
*/
