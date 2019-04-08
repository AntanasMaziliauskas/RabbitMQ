package person

import "github.com/AntanasMaziliauskas/RabbitMQ/types"

type PersonService interface {
	Init() error
	GetPerson(string) *types.Person
}
