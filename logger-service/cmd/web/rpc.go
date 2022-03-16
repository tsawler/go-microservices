package main

import (
	"context"
	"log"
	"time"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

// LogEntry is the type for all data stored in the logs collection. Note that we specify
// specific bson values, and we *must* include omitempty on ID, or newly inserted records will
// have an empty id!
type LogEntry struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Data      string    `bson:"data"`
	CreatedAt time.Time `bson:"created_at"`
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
		log.Println("Error writing log in rpc.go, LogInfo", err)
		return err
	}

	// resp is the message sent back to the RPC caller
	*resp = "Processed payload: " + payload.Name
	return nil
}
