package person

import (
	"context"
	"log"

	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataFromMongo struct {
	Name string
	Mgo  *mongo.Collection
}

func (d *DataFromMongo) Init() error {
	//Connect to DB
	if err := d.connectToDB(); err != nil {
		log.Fatal("Error while connecting to Database: ", err)
	}

	return nil
}

//GetOnePerson function looks for person and returns it if found
func (d *DataFromMongo) GetPerson(s string) *types.Person {
	result := &types.Person{}
	log.Println(s)
	if err := d.Mgo.FindOne(context.TODO(), bson.D{{"name", s}}).Decode(&result); err != nil {
		log.Println("Unable to locate given person")

		return &types.Person{}
	}

	return &types.Person{}

}

//connectTODB function connects to Mongo dabase.
func (d *DataFromMongo) connectToDB() error {

	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.200.244:27017"))
	if err != nil {
		panic(err)
	}
	//session.SetMode(mgo.Monotonic, true)
	d.Mgo = client.Database(d.Name).Collection("people")

	return nil
}
