package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/HugoW5/grpc-microservices/github.com/HugoW5/grpc-microservices/greetingpb"
	"github.com/HugoW5/grpc-microservices/github.com/HugoW5/grpc-microservices/userpb"
	"google.golang.org/grpc"
)

type server struct {
	greetingpb.UnimplementedGreetingServiceServer
	userClient userpb.UserServiceClient
}

func (s *server) SayHello(ctx context.Context, req *greetingpb.SayHelloRequest) (*greetingpb.SayHelloResponse, error) {
	userResp, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId})
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Hello, %s!", userResp.Name)
	return &greetingpb.SayHelloResponse{Message: message}, nil
}

func main() {
	conn, err := grpc.Dial("user-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	userClient := userpb.NewUserServiceClient(conn)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetingpb.RegisterGreetingServiceServer(s, &server{userClient: userClient})

	log.Println("Greeting Service running on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
