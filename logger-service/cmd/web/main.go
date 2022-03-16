package main

import (
	"context"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

var client *mongo.Client

const webPort = ":80"
const mongoURL = "mongodb://mongo:27017"

type Config struct {
	Session *scs.SessionManager
	Models  data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	client = mongoClient

	// we'll use this context to disconnect from mongo, since it needs one
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	// close connection to Mongo when application exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// set up app
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app := Config{
		Session: session,
		Models:  data.New(client),
	}

	// start webserver in its own GoRoutine
	go serve(app)

	// register the RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		return
	}

	// listen for RPC connections
	log.Println("Starting RPC Server on port 5001")
	listen, err := net.Listen("tcp", "0.0.0.0:5001")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	// this loop executes forever, waiting for connections
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		log.Println("Working...")
		go rpc.ServeConn(rpcConn)
	}
}

// serve starts the web server
func serve(app Config) {
	srv := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	fmt.Println("Starting logging web service on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connect opens a connection to mongo
func connectToMongo() (*mongo.Client, error) {
	// create connect options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect to the MongoDB and return Client instance
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		return nil, err
	}

	return c, nil
}
