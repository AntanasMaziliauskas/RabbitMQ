package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID         string
	Name       string
	Age        int
	Profession string
}

type PersonDB struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Age        int                `bson:"age"`
	Profession string             `bson:"profession"`
}

//Node structure holds values of Port, LastSeen, LastPing, IsOnline and Connection
type Node struct {
	//Port     string
	LastSeen time.Time
	//LastPing   time.Time //TODO: Ar tikrai jos reikia?
	IsOnline bool
	//Connection *grpc.ClientConn
	//Client     api.ServerClient
}
