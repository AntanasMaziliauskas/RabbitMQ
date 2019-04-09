package person

import (
	"context"
	"fmt"
	"log"

	"github.com/AntanasMaziliauskas/RabbitMQ/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	//log.Println(s)
	id, err := primitive.ObjectIDFromHex(s)
	//log.Println(id)
	if err != nil {
		log.Println("Invalid Object ID")

		return &types.Person{}
	}
	if err := d.Mgo.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&result); err != nil {
		log.Println("Unable to locate given person. Error: %s", err)

		return &types.Person{}
	}

	return result

}

func (d *DataFromMongo) UpsertPerson(in *types.Person) error {
	//var p Person

	p := primitive.D{
		{"$set", primitive.D{
			//{"_id", bson.ObjectIdHex(in.Id)},
			{"name", in.Name},
			{"age", in.Age},
			{"profession", in.Profession},
		}},
	}

	findOptions := options.Update()
	b := true
	findOptions.Upsert = &b
	id, err := primitive.ObjectIDFromHex(in.ID)
	/*if err != nil {
		log.Println("Object ID is invalid. Error: %s", err)
		return nil
	}*/
	//selector := bson.M{"_id": p.ID}
	//upsertdata := bson.M{"$set": p}
	_, err = d.Mgo.UpdateOne(context.TODO(), primitive.D{{"_id", id}}, p, findOptions)
	if err != nil {
		fmt.Println("Error while trying to upsert: ", err)

		return nil
	}
	//fmt.Println("Person successfully upserted")

	return nil
}

//GetOnePerson function looks for person and returns it if found
func (d *DataFromMongo) ListPersons() ([]types.Person, error) {
	var err error

	listOfData := []types.Person{}
	//Start
	//startTime := time.Now()

	cur, err := d.Mgo.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
	}
	//duration := time.Now().Sub(startTime)
	//log.Println(d.ID, ": Find funkcija: ", duration)
	for cur.Next(context.Background()) {

		elem := types.PersonDB{}

		//Start
		//startTime := time.Now()
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}
		//duration1 := time.Now().Sub(startTime)
		//log.Println(d.ID, ": Pries append: ", duration1)
		//list = append(list, &elem)
		listOfData = append(listOfData, types.Person{ID: elem.ID.Hex(),
			Name: elem.Name, Age: elem.Age, Profession: elem.Profession})

		//duration := time.Now().Sub(startTime)
		//log.Println(d.ID, ": vieno iraso apdorojimas Cursor: ", duration)
	}
	cur.Close(context.Background())

	//finish
	//duration := time.Now().Sub(startTime)
	//log.Println(d.ID, ": Kartu su Cursor paieška duombazėj užtruko: ", duration)

	if len(listOfData) < 1 {
		log.Println("This Node has no data.")
	}

	return listOfData, nil
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
