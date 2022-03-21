package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"log-service/logs"
	"net"
)

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
	log.Println("Log:", input.Name, input.Data)

	// TODO write the log

	res := &logs.LogResponse{Result: "logged!"}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", gRpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{})

	log.Printf("gRPC server started at port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
