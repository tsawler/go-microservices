package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

var infoLog *log.Logger
var client *mongo.Client

func main() {
	infoLogFile, err := os.OpenFile("./logs/logger/info_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening error log file: %v", err)
	}
	defer infoLogFile.Close()
	infoLog = log.New(infoLogFile, "INFO\t", log.Ldate|log.Ltime)

	mongoClient, err := connect()
	client = mongoClient

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	// close connection when func exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Println("Starting RPC Server on port 5001")
	err = rpc.Register(new(RPCServer))
	if err != nil {
		return
	}
	listen, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connect() (*mongo.Client, error) {
	// create connect options
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect to the MongoDB and return Client instance
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		return nil, err
	}

	return client, nil
}
