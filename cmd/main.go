package main

import (
	"log"
	"net"
	"net/http"
	"os"

	pb "service1/pb"
	"service1/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	myService := &service.MyService{}

	logger := log.New(os.Stdout, "Main: ", log.LstdFlags)

	grpcServer := grpc.NewServer()

	pb.RegisterMyServiceServer(grpcServer, myService)

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	logger.Println("Admin&User-Authentication-Server is running on :5050")
	go grpcServer.Serve(listener)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		logger.Fatalf("Failed to start health check server: %v", err)
	}
}
