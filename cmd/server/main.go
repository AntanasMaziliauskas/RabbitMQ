package main

import (
	"os"

	"github.com/AntanasMaziliauskas/RabbitMQ/server"
	"github.com/AntanasMaziliauskas/RabbitMQ/server/broker"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	service := server.Application{B: &broker.BrokerClient{}}
	service.Init()

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "getpersonnode",
			Usage: "Returns information about person from  Node that is specified",
			Action: func(c *cli.Context) error {
				return service.GetPerson(c)
			},
		},
		{
			Name:  "upsertperson",
			Usage: "Returns information about person from  Node that is specified",
			Action: func(c *cli.Context) error {
				return service.UpsertPerson(c)
			},
		},
		{
			Name:  "listpersonsbroadcast",
			Usage: "Returns information about person from  Node that is specified",
			Action: func(c *cli.Context) error {
				return service.ListPersonsBroadcast(c)
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "node",
			Value: "",
			Usage: "Name of the Node",
		},
		cli.StringFlag{
			Name:  "person",
			Value: "",
			Usage: "One or a list of names. For a list use ','. If you are inserting person use '.' to specify age and profession",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Starting application")
	}
}
