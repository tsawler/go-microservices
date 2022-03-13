package main

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

func (r *RPCServer) LogInfo(payload string, resp *string) error {
	infoLog.Println("Processed payload:", payload)
	*resp = "Processed payload: " + payload
	return nil
}
