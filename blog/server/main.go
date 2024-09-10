package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/inabinash/grpc/blog/proto"
	"go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("respones generated ", res)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	fmt.Printf("oid : %v ok %v", oid, ok)
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

func (*Server) ReadBlog(ctx context.Context, in *pb.BlogId) (*pb.Blog, error) {
	
	oid ,err := primitive.ObjectIDFromHex(in.Id);
	if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid Blog ID: %v", err))
    }
	result := &BlogItem{}

	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(
				codes.NotFound, fmt.Sprintf("Blog with this blog id not found %v", in.Id),
			)
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error finding blog: %v", err))
	}

	fmt.Println("got the result :", result)
	return DocumentToBlog(result), nil

}

func (*Server) ListBlog(in *empty.Empty,stream pb.BlogService_ListBlogServer) error {
	cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return status.Errorf(codes.Internal, fmt.Sprintf("Error retrieving blogs: %v", err))
    }
    defer cursor.Close(context.Background())

    // Iterate through the collection and stream each blog
    for cursor.Next(context.Background()) {
        result := &BlogItem{}
        err := cursor.Decode(result)
        if err != nil {
            return status.Errorf(codes.Internal, fmt.Sprintf("Error decoding blog: %v", err))
        }

        // Convert BlogItem to Blog proto message
        blog := DocumentToBlog(result)

        // Send each blog over the stream
        err = stream.Send(blog)
        if err != nil {
            return status.Errorf(codes.Internal, fmt.Sprintf("Error streaming blog: %v", err))
        }
    }

    if err := cursor.Err(); err != nil {
        return status.Errorf(codes.Internal, fmt.Sprintf("Cursor error: %v", err))
    }


	return nil;
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
