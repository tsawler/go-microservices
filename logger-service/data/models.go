package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// client is our mongo client that allows us to perform operations on the mongo database
var client *mongo.Client

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *mongo.Client) Models {
	client = dbPool

	return Models{
		LogEntry: LogEntry{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	LogEntry LogEntry
}

// LogEntry is the type for all data stored in the logs collection. Note that we specify
// specific bson values, and we *must* include omitempty on ID, or newly inserted records will
// have an empty id!
type LogEntry struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Data      string    `bson:"data"`
	CreatedAt time.Time `bson:"created_at"`
}

// Insert puts a document in the logs collection
func (l *LogEntry) Insert(entry LogEntry) (string, error) {
	collection := client.Database("logs").Collection("logs")

	result, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry:", err)
		return "", err
	}

	return result.InsertedID.(string), nil
}

// All returns all documents in the log collection
func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all documents ERROR:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error scanning log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
