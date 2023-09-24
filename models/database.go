package models

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB         *mongo.Client
	Collection map[string]*mongo.Collection
)

func ConnectDatabase() {
	var err error
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Connection Failed to Database")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Unable to connect to Database")
		log.Fatal(err)
	}

	log.Println("⛁ Connected to Mongo Database!")
	collections := loadCollection(client)
	DB = client
	Collection = collections
}

func loadCollection(mongoConn *mongo.Client) map[string]*mongo.Collection {
	collections := make(map[string]*mongo.Collection, 4)
	collections["patients"] = colHelper(mongoConn, "patients")
	collections["patient_form"] = colHelper(mongoConn, "patient_form")
	collections["files"] = colHelper(mongoConn, "files")
	collections["forms"] = colHelper(mongoConn, "forms")

	return collections
}

func colHelper(db *mongo.Client, collectionName string) *mongo.Collection {
	return db.Database("patient").Collection(collectionName)
}
