package server

import (
	"log"

	"github.com/AntanasMaziliauskas/RabbitMQ/server/broker"
	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"github.com/urfave/cli"
)

type Application struct {
	B broker.BrokerService
}

func (a *Application) Init() error {
	a.B.Init()

	return nil
}

func (a *Application) GetPerson(c *cli.Context) error {
	if c.GlobalString("node") != "" && c.GlobalString("name") != "" {
		result, _ := a.B.GetPersonNode(c.GlobalString("node"), "getPerson", types.Person{Name: c.GlobalString("name")})

		log.Println(result)
	}
	a.B.Stop()
	return nil
}

func (a *Application) UpsertPerson(c *cli.Context) error {
	if c.GlobalString("node") != "" && c.GlobalString("name") != "" {
		result, _ := a.B.UpsertPerson(c.GlobalString("node"), "getPerson", types.Person{Name: c.GlobalString("name")})

		log.Println(result)
	}
	a.B.Stop()
	return nil
}
