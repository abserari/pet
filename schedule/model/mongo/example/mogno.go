package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// You will be using this Trainer type later in the program
type Schedule struct {
	ID      uint64
	AdminID uint64
	PetID   uint64
	Date    time.Time
	Title   string
	Note    string
	Shop    uint64
}

func main() {
	// Rest of the code will go here
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://root:yhyddr119216@cluster0.ynf1u.mongodb.net/test?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("schedule")
	moliKanbing := Schedule{100, 1000, 10000, time.Now(), "例行检查", "无笔记", 0}

	insertResult, err := collection.InsertOne(context.TODO(), moliKanbing)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Insert a single document", insertResult.InsertedID)

	filter := bson.D{{"id", 100}}

	update := bson.D{
		{
			"$set", bson.D{
				{"note", "有笔记"},
			}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// create a value into which the result can be decoded
	var result Schedule

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
