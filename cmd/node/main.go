package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntanasMaziliauskas/RabbitMQ/node"
	"github.com/AntanasMaziliauskas/RabbitMQ/node/person"
)

func main() {
	name := flag.String("name", "Node01", "Name of the Node")
	flag.Parse()

	app := node.Application{
		Name: *name,
		Person: &person.DataFromMongo{
			Name: *name,
		},
	}

	//Pradedam klausytis
	app.Init()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGKILL)

	<-stop
}
