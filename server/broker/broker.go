package broker

import (
	"bytes"
	"encoding/gob"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"github.com/streadway/amqp"
)

type BrokerService interface {
	Init() error
	Stop() error
	GetPersonNode(string, string, types.Person) (types.Person, error)
	UpsertPerson(string, string, types.Person) error
	ListPersonsBroadcast() ([]types.Person, error)
	ListenForGreeting()
}

type BrokerClient struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	Reply <-chan amqp.Delivery
	Ping  <-chan amqp.Delivery
	wg    *sync.WaitGroup
	Nodes map[string]*types.Node
}

/*

func (b *BrokerClinet) Call(node. command string, payload interface{}) (interface{}, error) {

	corrID := genNewCorrID()

//
	xxx := struct{
		chan,
		created time.Time,
		count int,
	}
//



	reply := make(chan, amqp.Delivery, 1)
	timeout := time.NewTicker(1s)
	b.callbacks[corrID] := reply

	gob := GobMarshal(payload)
	b.ch.Publish(
		...
		body: gob,
		Header.CorrID: corrID
		...
	)


	// veliau
	for count > 0 || timeout {
		packet := <- b.callbacks[corrID]
		x := gob.Unmarshal(packet.Body)
		data[nodeID] = x

		coun--
	}

	select {
	case packet := <- b.callbacks[corrID]:
		x := gob.Unmarshal(packet.Body)
		return x, nil
	case <- timeout.C:
		close(b.callbacks[corr])
		timeout.Stop()
		delete(b.callbacks[corr])
		return nil, fmt.Error("timeouted")
	}



	return data, nil
}



func (b *BrokerClient) Listen() {
		for d := range b.Reply {
			log.Printf("Server %s received reply.")
			if b.callbacks[d.CorrID] {
				b.callbacks[d.CorrID] <- d
			}
		}
}
*/

func (b *BrokerClient) Init() error {
	var err error

	b.Nodes = make(map[string]*types.Node)
	rand.Seed(time.Now().UTC().UnixNano())
	b.Conn = &amqp.Connection{}
	b.Ch = &amqp.Channel{}
	b.Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Failed to connect to RabbitMQ")

	b.Ch, err = b.Conn.Channel()
	handleError(err, "Failed to open a channel")

	err = b.Ch.ExchangeDeclare(
		"cmd",   // name
		"topic", // type
		false,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	handleError(err, "Failed to declare an exchange")

	//response queue
	q, err := b.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	handleError(err, "Failed to declare a queue")

	//ping queue
	p, err := b.Ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	handleError(err, "Failed to declare a queue")

	//Bindinam Ping queue
	err = b.Ch.QueueBind(
		p.Name, // queue name
		"ping", // routing key
		"cmd",  // exchange
		false,
		nil)
	handleError(err, "Failed to bind a queue")

	//Bindinam response queue
	err = b.Ch.QueueBind(
		q.Name,     // queue name
		"response", // routing key
		"cmd",      // exchange
		false,
		nil)
	handleError(err, "Failed to bind a queue")

	b.Ping, err = b.Ch.Consume(
		p.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")

	b.Reply, err = b.Ch.Consume(
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

func (b *BrokerClient) ListenForGreeting() {
	//b.wg.Add(1)
	go func() {
		//defer b.wg.Done()
		for d := range b.Ping {
			log.Printf("Node %s connected.", string(d.Body))
			b.Nodes[string(d.Body)] = &types.Node{
				LastSeen: time.Now(),
				IsOnline: true,
			}
			d.Ack(false)
		}
	}()
}

func (b *BrokerClient) GetPersonNode(to string, action string, person types.Person) (types.Person, error) {
	res := types.Person{}

	corrId := randomString(32)
	per, _ := GobMarshal(person)

	log.Println(to)
	//to = "nodes." + to
	err := b.Ch.Publish(
		"cmd", // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"Node": to, "action": "getPerson"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			//ReplyTo:       "response",  //q.Name
			Body: []byte(per), //Convert
		})

	handleError(err, "Failed to publish a message")

	for d := range b.Reply {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			//res.Profession = string(d.Body)
			break
		}
	}

	return res, nil
}

func (b *BrokerClient) UpsertPerson(to string, action string, person types.Person) error {
	//res := types.Person{}

	//	corrId := randomString(32)
	per, _ := GobMarshal(person)

	err := b.Ch.Publish(
		"cmd", // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "upsertPerson"},
			ContentType:   "text/plain",
			CorrelationId: "corrId",
			ReplyTo:       "response",
			Body:          []byte(per), //Convert
		})
	handleError(err, "Failed to publish a message")

	/*for d := range b.Reply {
		if corrId == d.CorrelationId {
			GobUnmarshal(d.Body, &res)
			//res.Profession = string(d.Body)
			break
		}
	}*/

	return nil
}

//TODO: FANOUT
func (b *BrokerClient) ListPersonsBroadcast() ([]types.Person, error) {
	res := []types.Person{}
	result := []types.Person{}

	corrId := randomString(32)
	//per, _ := GobMarshal(person)

	err := b.Ch.Publish(
		"cmd",       // exchange
		"broadcast", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			Headers:       map[string]interface{}{"action": "listPersons"},
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       "response",
			//Body:          []byte(per), //Convert
		})
	handleError(err, "Failed to publish a message")

	//Kaip cia viskas atrodys?
	i := 1
	for d := range b.Reply {
		if corrId == d.CorrelationId {
			i++
			GobUnmarshal(d.Body, &res)
			result = append(result, res...)
			//res.Profession = string(d.Body)
			if i > len(b.Nodes) {
				break
			}
		}
	}
	//log.Println(result)
	return result, nil
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
