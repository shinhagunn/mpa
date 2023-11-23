package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDBUser     = "mongo"
	MongoDBPassword = "mongo"
	MongoDBHost     = "localhost"
	MongoDBPort     = "27017"
)

type Person struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
}

type Address struct {
	PersonID primitive.ObjectID `bson:"person_id"`
	City     string             `bson:"city"`
}

// func ApplyFilters(filters ...bson.E) bson.D {
// 	var result []bson.E
// 	result = append(result, filters...)
// 	return result
// }

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up MongoDB connection options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", MongoDBUser, MongoDBPassword, MongoDBHost, MongoDBPort))

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Disconnect from MongoDB when the application exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Access a MongoDB collection
	collectionPeople := client.Database("testdb").Collection("people")
	_ = client.Database("testdb").Collection("addresses")

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "addresses"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "person_id"},
			{Key: "as", Value: "addresses"},
		}},
	}

	// Add additional stages or filters if needed
	pipeline := mongo.Pipeline{lookupStage}

	cursor, err := collectionPeople.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var peopleWithAddresses []bson.M
	if err := cursor.All(ctx, &peopleWithAddresses); err != nil {
		log.Fatal(err)
	}

	fmt.Println("People with addresses:")
	for _, result := range peopleWithAddresses {
		fmt.Printf("%v\n", result)
	}
}
