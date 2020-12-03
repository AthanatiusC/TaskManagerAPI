package models

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	db = client.Database("TaskManager")
}

func GetDB(whatDB string) *mongo.Database {
	if whatDB == "main" {
		return db
	} else {
		return db
	}
	return db
}
