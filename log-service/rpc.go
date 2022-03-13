package main

import (
	"context"
	"log"
	"time"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

type LogEntry struct {
	Data      string
	CreatedAt time.Time
}

func (r *RPCServer) LogInfo(payload string, resp *string) error {
	infoLog.Println("Processed payload:", payload)

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Data:      payload,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	*resp = "Processed payload: " + payload
	return nil
}
