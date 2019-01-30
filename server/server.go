package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	metricsInterface = "0.0.0.0:5117"
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
	grpcInterface := flag.Arg(0)

	go func() {
		l, err := net.Listen("tcp", metricsInterface)
		if err != nil {
			log.Fatalln(err)
		}
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		s := &http.Server{
			Handler: mux,
		}
		log.Printf("Serving HTTP /metrics on %s\n", l.Addr())
		log.Fatalln(s.Serve(l))
	}()

	func() {
		l, err := net.Listen("tcp", grpcInterface)
		if err != nil {
			log.Fatalln(err)
		}
		server := grpc.NewServer(
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		)
		grpc_prometheus.Register(server)

		RegisterGreeterServer(server, &Service{})

		healthsrv := health.NewServer()
		healthsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
		healthpb.RegisterHealthServer(server, healthsrv)

		log.Printf("Serving gRPC on %s\n", l.Addr())
		log.Fatalln(server.Serve(l))
	}()
}
