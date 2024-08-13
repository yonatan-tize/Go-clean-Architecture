package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectDB establishes a connection to the MongoDB database using the provided URI.
// It returns a pointer to the mongo.Database object representing the connected database.
// The URI parameter specifies the connection string for the MongoDB server.
// The function logs a fatal error and terminates the program if it fails to connect to the database.
// The function also sets the global variable Client to the connected client for future use.
// The function assumes that the database name is "User_Task_Manager".
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


// DisconnectDB disconnects the database connection.
// It calls the Disconnect method on the Client object and logs an error if it occurs.
func DisconnectDB(){
	if err := Client.Disconnect(context.TODO()); err != nil{
		log.Fatal("failed to disconnect database")
	}
	
}