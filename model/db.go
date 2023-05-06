package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func InitMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://11.11.9.202:27017"))
	Db = MongoClient.Database("energy")
	if err != nil {
		log.Panicf("MongoDB connection error: %s", err)
	}
}
