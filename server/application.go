package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/AntanasMaziliauskas/RabbitMQ/server/broker"
	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"github.com/urfave/cli"
)

type Application struct {
	B  broker.BrokerService
	wg *sync.WaitGroup
}

func (a *Application) Init() error {
	a.B.Init()

	return nil
}

func (a *Application) GetPerson(c *cli.Context) error {
	if c.GlobalString("node") != "" && c.GlobalString("person") != "" {
		//person := parsePerson(c.GlobalString("person"))
		result, _ := a.B.GetPersonNode(c.GlobalString("node"), "getPerson", types.Person{ID: c.GlobalString("person")})

		log.Println(result)
	}
	a.B.Stop()
	return nil
}

func (a *Application) UpsertPerson(c *cli.Context) error {
	if c.GlobalString("node") != "" && c.GlobalString("person") != "" {
		person := parsePerson(c.GlobalString("person"))
		a.B.UpsertPerson(c.GlobalString("node"), "upsertPerson", person)

		//log.Println(result)
	}
	//log.Println("Pries STOP")
	a.B.Stop()
	return nil
}

func (a *Application) ListPersonsBroadcast(c *cli.Context) error {

	result, _ := a.B.ListPersonsBroadcast()

	log.Println(result)

	a.B.Stop()
	return nil
}

//parsePerson function parses string into Person structure
func parsePerson(s string) types.Person {
	var (
		person types.Person
		age    int
		err    error
	)

	slice := strings.Split(s, ".")
	if len(slice) > 1 {
		//i, err := strconv.Atoi("-42")
		if age, err = strconv.Atoi(slice[2]); err != nil {
			fmt.Println("Error while concerting string into int: ", err)
		}
		person = types.Person{
			ID:         slice[0],
			Name:       slice[1],
			Age:        age,
			Profession: slice[3],
		}
		return person
	}
	person = types.Person{
		ID: slice[0],
	}
	return person
}

//parsePersons function parses string into []Person structure
func parsePersons(s string) []types.Person {
	var (
		list []types.Person
	)

	persons := strings.Split(s, ",")
	for _, k := range persons {
		one := parsePerson(k)
		list = append(list, one)
	}
	return list
}
