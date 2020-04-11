package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/chiyukuan/golab/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetN(req *greetpb.GreetRequest, stream greetpb.GreetService_GreetNServer) error {

	firstName := req.GetGreeting().GetFirstName()
	for ii := 0; ii < 10; ii++ {
		result := "Hello " + firstName + " @ " + strconv.Itoa(ii)
		res := &greetpb.GreetResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LGreet(stream greetpb.GreetService_LGreetServer) error {

	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&greetpb.GreetResponse{Result: result})
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
	return nil
}

func (*server) LGreetN(stream greetpb.GreetService_LGreetNServer) error {

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		stream.Send(&greetpb.GreetResponse{Result: result})
	}
	return nil
}

func main() {
	fmt.Println("vim-go")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, errCert := credentials.NewServerTLSFromFile(certFile, keyFile)
	if errCert != nil {
		log.Fatalf("Error while loading certificate: %v", errCert)
	}
	s := grpc.NewServer(grpc.Creds(creds))

	greetpb.RegisterGreetServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
