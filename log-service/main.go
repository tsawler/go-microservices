package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func main() {
	//go listenRPC()
	log.Println("Starting RPC Server on port 5001")
	err := rpc.Register(new(RPCServer))
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
