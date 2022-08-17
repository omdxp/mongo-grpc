package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/Omar-Belghaouti/mongo-grpc/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var collection *mongo.Collection

type server struct {
	pb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create mongodb client
	log.Print("connecting to mongodb")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err.Error())
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		log.Print("starting server")
		if err = s.Serve(lis); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch

	log.Print("stopping server")
	s.Stop()

	log.Print("closing listener")
	lis.Close()

	log.Print("closing mongodb connection")
	if err = client.Disconnect(ctx); err != nil {
		panic(err)
	}

	log.Print("bye bye")

}
