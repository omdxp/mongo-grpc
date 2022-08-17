package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Omar-Belghaouti/mongo-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedBlogServiceServer
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	log.Print("bye bye")

}
