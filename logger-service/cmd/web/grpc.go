package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"
)

// LogServer is type used for writing to the log via gRPC. Note that we embed the
// data.Models type, so we have access to Mongo.
type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

// WriteLog writes the log after receiving a call from a gRPC client. This function
// must exist, and is defined in logs/logs.proto, in the "service LogService" bit
// at the end of the file.
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
	log.Println("Log:", input.Name, input.Data)

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed:"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "logged!"}

	return res, nil
}

// gRPCListen starts the gRPC server
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// register the service, handing it models (so we can write to the database)
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC server started at port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
