package database

import (
	"fmt"
	"log"
	"os"

	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DBInstance() *mongo.Client {
	err:=gotenv.Load(".env")

	if err!=nil {
		log.Println("Warning: unable to fund .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")

	if MongoDb == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	
	fmt.Println("MONGODB_URI", MongoDb)

	clientOptions := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}

	return client

}

var Client *mongo.Client = DBInstance()

func OpenCollection( collectionName string) *mongo.Collection {
	err:=gotenv.Load(".env")

	if err!=nil {
		log.Println("Warning: unable to fund .env file")
	}


	dbName := os.Getenv("DATABASE_NAME")

	fmt.Println("dbName", dbName)

	collection := Client.Database(dbName).Collection(collectionName)

	if err != nil {
		log.Fatal(err)
	}

	return collection
}