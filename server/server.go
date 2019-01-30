package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeeraw/greeter/server/controllers"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	metricsInterface = "0.0.0.0:5117"
)

var (
	errTimeout = status.Error(codes.DeadlineExceeded, "request took too long")
)

// Service represents the greeter service
type Service struct{}

// Hello responds to a greeting
func (s *Service) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	controller := &controllers.Greetings{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "metadata missing")
	}
	jwts := md.Get("jwt")
	if len(jwts) < 1 {
		return nil, status.Error(codes.Unauthenticated, "jwt missing")
	}
	greeting, err := controller.Hello(ctx, req.Name)
	if err != nil {
		switch e := err.(type) {
		default:
			return nil, status.Error(codes.Internal, e.Error())
		}
	}
	return &HelloResponse{Greeting: greeting}, nil
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
