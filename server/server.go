package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeeraw/greeter/server/controllers"

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

var (
	errTimeout = status.Error(codes.DeadlineExceeded, "request took too long")
)

// Service represents the greeter service
type Service struct{}

// Hello responds to a greeting
func (s *Service) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	controller := &controllers.Greetings{}
	res := controller.Hello(req.Name)
	defer res.Close()

	select {
	case err := <-res.ErrC:
		return nil, status.Error(codes.Internal, err.Error())
	case greeting := <-res.ResC:
		return &HelloResponse{Greeting: greeting}, nil
	case <-ctx.Done():
		return nil, errTimeout
	}
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
