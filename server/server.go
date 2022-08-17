package main

import (
	"log"
	"net"

	"github.com/Omar-Belghaouti/mongo-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedBlogServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	if err = s.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
