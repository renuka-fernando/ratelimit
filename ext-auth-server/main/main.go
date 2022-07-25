package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"

	authsample "extauth/sample"
)

func main() {
	port := flag.Int("port", 9001, "gRPC port")
	data := flag.String("users", "users.json", "users file")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", *port, err)
	}

	users, err := authsample.LoadUsers(*data)
	if err != nil {
		log.Fatalf("failed to load user data:%s %v", *data, err)
	}
	gs := grpc.NewServer()

	envoy_service_auth_v3.RegisterAuthorizationServer(gs, authsample.New(users))

	log.Printf("starting gRPC server on: %d\n", *port)

	gs.Serve(lis)
}
