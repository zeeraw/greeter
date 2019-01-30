package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Service represents the greeter service
type Service struct{}

// Hello responds to a greeting
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

	healthsrv := health.NewServer()
	healthsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthsrv)

	log.Printf("Serving GRPC on %s\n", listener.Addr())
	log.Fatalln(server.Serve(listener))
}
