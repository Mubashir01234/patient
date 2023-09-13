package models

import (
	"context"
	"log"
	"patient/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB         *mongo.Client
	Collection map[string]*mongo.Collection
)

func ConnectDatabase() {
	var err error
	clientOptions := options.Client().ApplyURI(config.Cfg.MongoURL)
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
	collections := make(map[string]*mongo.Collection, 2)
	collections["users"] = colHelper(mongoConn, "users")
	collections["patient_form"] = colHelper(mongoConn, "patient_form")
	return collections
}

func colHelper(db *mongo.Client, collectionName string) *mongo.Collection {
	return db.Database("patient").Collection(collectionName)
}
