package main

import (
	"log"

	"github.com/Omar-Belghaouti/mongo-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}

	c := pb.NewBlogServiceClient(cc)
	_ = c
}
