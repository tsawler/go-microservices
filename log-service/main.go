package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

var infoLog *log.Logger

func main() {
	infoLogFile, err := os.OpenFile("./logs/logger/info_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening error log file: %v", err)
	}
	defer infoLogFile.Close()
	infoLog = log.New(infoLogFile, "INFO\t", log.Ldate|log.Ltime)

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
