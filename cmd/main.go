package main

import (
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	pb "service1/pb"
	"service1/pkg/service"
)

func main() {
	myService := &service.MyService{}

	grpcServer := grpc.NewServer()

	pb.RegisterMyServiceServer(grpcServer, myService)

	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Admin&User-Authentication-Server is running on 5050")
	go grpcServer.Serve(listener)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start health check server: %v", err)
	}
}
