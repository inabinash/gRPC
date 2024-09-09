package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/inabinash/grpc/blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlogItem struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Author  string             `bson:"author"`
	Content string             `bson:"content"`
}

func DocumentToBlog(data *BlogItem) *pb.Blog {
	return &pb.Blog{
		Id:      data.Id.Hex(),
		Author:  data.Author,
		Name:    data.Name,
		Content: data.Content,
	}
}

type Server struct {
	pb.UnimplementedBlogServiceServer
}

var (
	port       = flag.Int("port", 50051, "The server port")
	collection *mongo.Collection
)

func (*Server) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	fmt.Println("create blog function invoked")

	data := BlogItem{
		Name:    in.Name,
		Author:  in.Author,
		Content: in.Content,
	}
	fmt.Println("-----------------data------------", data)
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"Cannot convert to OID",
		)
	}

	return &pb.BlogId{
		Id: oid.Hex(),
	}, nil
}
func connectDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))

	return client, err
}
func main() {
	// connect to database
	//create a collection
	client, err := connectDB()

	if err != nil {
		log.Fatal("error in creating the client with the given uri :", err)
	}
	collection = client.Database("blogdb").Collection("blogs")
	fmt.Printf("collection is %+v\n", *collection)

	//  write the create

	// Listen to the requests

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		log.Fatal("Failed to open a network with the provided address .", err)
	}
	fmt.Printf("listning object ... %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve")
	}

	// and read operations

}
