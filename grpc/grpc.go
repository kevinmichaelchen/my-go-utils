package grpc

import (
	"time"
	"log"
	"google.golang.org/grpc"
)

// InitGrpcConn creates a gRPC connection.
//
// Example usage:
//
//   conn := utils.InitGrpcConn("localhost:50051", 3, time.Second*5)
func InitGrpcConn(address string, numRetries int, sleepDuration time.Duration) *grpc.ClientConn {
	var err error
	var conn *grpc.ClientConn
	for i := 0; i < numRetries; i++ {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err == nil {
			log.Println("Grpc connection initialized...")
			break
		}
		log.Println("Grpc connection failed to initialize... Sleeping...")
		time.Sleep(sleepDuration)
	}
	if err != nil {
		log.Fatal(err)
	}
	return conn
}