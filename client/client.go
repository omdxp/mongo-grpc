package main

import (
	"context"
	"io"
	"log"

	"github.com/Omar-Belghaouti/mongo-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cc.Close()

	c := pb.NewBlogServiceClient(cc)

	// create a blog
	// createBlog(c)

	// read a blog
	// readBlog(c)

	// update a blog
	// updateBlog(c)

	// delete a blog
	// deleteBlog(c)

	// list blogs
	listBlogs(c)
}

func createBlog(c pb.BlogServiceClient) {
	req := &pb.CreateBlogRequest{
		Blog: &pb.Blog{
			AuthorId: "Omar",
			Title:    "My First Blog",
			Content:  "Content of the first blog",
		},
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Print(s.Code(), ": ", s.Message())
		} else {
			log.Fatal(err.Error())
		}
	}
	log.Printf("blog created: %v", res.GetBlog())
}

func readBlog(c pb.BlogServiceClient) {
	req := &pb.ReadBlogRequest{
		BlogId: "62fcaacf410e7788bd475335",
	}

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Print(s.Code(), ": ", s.Message())
		} else {
			log.Fatal(err.Error())
		}
	}

	log.Printf("blog found: %v", res.GetBlog())
}

func updateBlog(c pb.BlogServiceClient) {
	req := &pb.UpdateBlogRequest{
		Blog: &pb.Blog{
			Id:       "62fcaacf410e7788bd475335",
			AuthorId: "Ramo",
			Title:    "New Title",
			Content:  "New content",
		},
	}
	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Print(s.Code(), ": ", s.Message())
		} else {
			log.Fatal(err.Error())
		}
	}

	log.Printf("blog updated: %v", res.GetBlog())
}

func deleteBlog(c pb.BlogServiceClient) {
	req := &pb.DeleteBlogRequest{
		BlogId: "62fcaacf410e7788bd475335",
	}
	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Print(s.Code(), ": ", s.Message())
		} else {
			log.Fatal(err.Error())
		}
	}

	log.Printf("blog deleted: %s", res.GetBlogId())
}

func listBlogs(c pb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &pb.ListBlogRequest{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				log.Print(s.Code(), ": ", s.Message())
			} else {
				log.Fatal(err.Error())
			}
		}
		log.Printf("blog received: %v", res.GetBlog())
	}
}
