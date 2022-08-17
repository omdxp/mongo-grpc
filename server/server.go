package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/Omar-Belghaouti/mongo-grpc/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

// create a blog to mongodb
func (s *server) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	log.Print("CreateBlog invoked")
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %s", err.Error())
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "cannot convert to object id")
	}

	return &pb.CreateBlogResponse{
		Blog: &pb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

// read blog from mongodb
func (s *server) ReadBlog(ctx context.Context, req *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	log.Print("ReadBlog invoked")
	blodId := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blodId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot parse id")
	}

	var res *blogItem
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "document not found")
		}
		return nil, status.Errorf(codes.Internal, "internal err: %s", err.Error())
	}

	return &pb.ReadBlogResponse{
		Blog: &pb.Blog{
			Id:       res.ID.Hex(),
			AuthorId: res.AuthorID,
			Title:    res.Title,
			Content:  res.Content,
		},
	}, nil
}

// update blog in mongodb
func (s *server) UpdateBlog(ctx context.Context, req *pb.UpdateBlogRequest) (*pb.UpdateBlogResponse, error) {
	log.Print("UpdateBlog invoked")
	blog := req.GetBlog()

	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot parse id")
	}

	// find document by id
	_, err = s.ReadBlog(ctx, &pb.ReadBlogRequest{
		BlogId: oid.Hex(),
	})
	if err != nil {
		return nil, err
	}

	// update document
	data := &blogItem{
		ID:       oid,
		AuthorID: blog.AuthorId,
		Content:  blog.Content,
		Title:    blog.Title,
	}

	_, err = collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: oid}}, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal err: %s", err.Error())
	}

	return &pb.UpdateBlogResponse{
		Blog: &pb.Blog{
			Id:       oid.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *pb.DeleteBlogRequest) (*pb.DeleteBlogResponse, error) {
	log.Print("DeleteBlog invoked")
	blogId := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "cannot parse id")
	}

	// delete document
	res, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal err: %s", err.Error())
	}

	if res.DeletedCount == 0 {
		return nil, status.Error(codes.NotFound, "document not found")
	}

	return &pb.DeleteBlogResponse{
		BlogId: blogId,
	}, nil
}

func (s *server) ListBlog(req *pb.ListBlogRequest, stream pb.BlogService_ListBlogServer) error {
	log.Print("ListBlog invoked")

	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, "internal err: %s", err.Error())
	}
	defer cur.Close(context.Background())

	// iterate over data
	var data *blogItem
	for cur.Next(context.Background()) {
		err := cur.Decode(&data)
		if err != nil {
			return status.Errorf(codes.Internal, "internal err: %s", err.Error())
		}
		stream.Send(&pb.ListBlogResponse{
			Blog: &pb.Blog{
				Id:       data.ID.Hex(),
				AuthorId: data.AuthorID,
				Title:    data.Title,
				Content:  data.Content,
			},
		})
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(codes.Internal, "internal err: %s", err.Error())
	}
	return nil
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
