package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/chiyukuan/golab/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("vim-go")

	certFile := "ssl/ca.crt"
	creds, errCert := credentials.NewClientTLSFromFile(certFile, "")
	if errCert != nil {
		log.Fatalf("error while loading CA: %v", errCert)
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	greetClt := greetpb.NewGreetServiceClient(conn)
	//doUnary(greetClt)
	//doServerStreaming(greetClt)
	//doClientStreaming(greetClt)
	doBiDirStreaming(greetClt)
}

func doUnary(greetClt greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ray",
			LastName:  "Kuan",
		},
	}
	res, err := greetClt.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greeting RPC, %v", err)
	}
	log.Printf("Response from Greet: %v\n", res.GetResult())
}

func doServerStreaming(greetClt greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ray",
			LastName:  "Kuan",
		},
	}
	res, err := greetClt.GreetN(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetN RPC: %v", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			log.Printf("End of Streaming %v", msg)
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetN: %v\n", msg.GetResult())
	}
}

func doClientStreaming(greetClt greetpb.GreetServiceClient) {

	requests := []*greetpb.GreetRequest{
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Ray", LastName: "Kuan"}},
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Amy", LastName: "Kuan"}},
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Michelle", LastName: "Kuan"}},
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Constance", LastName: "Kuan"}},
	}

	stream, err := greetClt.LGreet(context.Background())

	if err != nil {
		log.Fatalf("")
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response message %v\n", err)
	}
	fmt.Printf("LGreet response: %v\n", res.GetResult())
}

func doBiDirStreaming(greetClt greetpb.GreetServiceClient) {

	requests := []*greetpb.GreetRequest{
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Ray", LastName: "Kuan"}},
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Amy", LastName: "Kuan"}},
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Michelle", LastName: "Kuan"}},
		&greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Constance", LastName: "Kuan"}},
	}

	stream, err := greetClt.LGreetN(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			stream.Send(req)
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Printf("End of Streaming %v", res)
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
			}
			log.Printf("Response from GreetN: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
