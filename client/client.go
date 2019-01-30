package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(flag.Arg(0), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Dialing GRPC on %s\n", conn.Target())
	client := NewGreeterClient(conn)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}))
	res, err := client.Hello(ctx, &HelloRequest{
		Name: flag.Arg(1),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Greeting)
}
