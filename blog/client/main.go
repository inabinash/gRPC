package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/inabinash/grpc/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addrs = "127.0.0.1:3000"

func createBlog(c pb.BlogServiceClient) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	req := &pb.Blog{
		Name:    "first blog",
		Author:  "Abinash",
		Content: "surviving in samsung",
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatal(" error in creating the blog ", err)
		return

	}
	fmt.Printf("Genrated a blog with blog id %v\n", res.Id)
}
func main() {
	conn, err := grpc.NewClient(addrs, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()
	c := pb.NewBlogServiceClient(conn)
	createBlog(c);
}
