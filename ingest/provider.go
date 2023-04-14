package ingest

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func InsertRecord(uri string) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// insert data
	collection := client.Database("comet").Collection("chain")
	if _, err = collection.InsertOne(ctx, bson.D{{"chain_id", "chain1"}}); err != nil {
		log.Fatal(err)
	}

	// show a list of databases
	databases, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Databases: %v", databases)

	// show inserted data
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var res []bson.D
	if err = cursor.All(ctx, &res); err != nil {
		log.Fatal(err)
	}
	log.Printf("Data: %v", res)
}
