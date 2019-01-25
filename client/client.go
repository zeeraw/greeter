package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(flag.Arg(0), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Dialing GRPC on %s\n", conn.Target())
	client := NewGreeterClient(conn)
	res, err := client.Hello(context.Background(), &HelloRequest{
		Name: flag.Arg(1),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Greeting)
}
