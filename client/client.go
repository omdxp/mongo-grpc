package main

import (
	"context"
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
	defer cc.Close()

	c := pb.NewBlogServiceClient(cc)

	// create a blog
	req := &pb.CreateBlogRequest{
		Blog: &pb.Blog{
			AuthorId: "Omar",
			Title:    "My First Blog",
			Content:  "Content of the first blog",
		},
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("blog created: %v", res.GetBlog())
}
