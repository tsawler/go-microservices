package main

import (
	"log"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

func (r *RPCServer) LogInfo(payload string, resp *string) error {
	log.Println("Processed rpc call")
	return nil
}
