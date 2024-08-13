package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB(uri string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("failed to connect")
	}
	
	Client = client
	return client.Database("User_Task_Manager") // name of the database
}

func DisconnectDB(){
	if err := Client.Disconnect(context.TODO()); err != nil{
		log.Fatal("failed to disconnect database")
	}
	
}