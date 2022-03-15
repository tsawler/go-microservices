package main

import (
	"context"
	"log"
	"time"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

// LogEntry is the type for writing data to mongo
type LogEntry struct {
	Name      string
	Data      string
	CreatedAt time.Time
}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo writes our payload to mongo
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	infoLog.Println("Processed payload:", payload)

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// resp is the message sent back to the RPC caller
	*resp = "Processed payload: " + payload.Name
	return nil
}
