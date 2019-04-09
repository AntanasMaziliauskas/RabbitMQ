package types

import "go.mongodb.org/mongo-driver/bson/primitive"

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
