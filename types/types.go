package types

type Person struct {
	Name       string `bson:"name"`
	Age        int    `bson:"age"`
	Profession string `bson:"profession"`
}
