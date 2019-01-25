package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Service struct{}

func (s *Service) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{
		Greeting: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}
	server := grpc.NewServer()
	RegisterGreeterServer(server, &Service{})
	log.Printf("Serving GRPC on %s\n", listener.Addr())
	log.Fatalln(server.Serve(listener))
}
