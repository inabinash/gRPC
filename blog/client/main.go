package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/inabinash/grpc/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var addrs = "127.0.0.1:3000"

func doCreate(c pb.BlogServiceClient) (string, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	req := &pb.Blog{
		Name:    "third blog",
		Author:  "abhishek",
		Content: "whatever grpc",
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatal(" error in creating the blog ", err)
		return "", err

	}
	fmt.Printf("Genrated a blog with blog id %v\n", res.Id)
	return res.Id, nil
}

func doRead(c pb.BlogServiceClient, id string) {
	req := &pb.BlogId{
		Id: id,
	}

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatal(" error in reading the blog ", err)
	}
	fmt.Printf("Genrated a blog with blog id %v\n name %v\n author %v\n content %v\n ", res.Id, res.Name, res.Author, res.Content)
}

func doList(c pb.BlogServiceClient) {
	req := &emptypb.Empty{}
	stream ,err:= c.ListBlog(context.Background(), req);
	if err != nil {
		log.Fatal(" error in listing the blog ", err);
	}
	for {
		blog , err:= stream.Recv();
		if err != nil {
			if err == io.EOF {
				log.Fatal("End of stream :", err);
				return;
			}
			log.Fatal("Unexpected error ", err);

		}
		fmt.Printf("Received blog with id %v name %v \n", blog.Id, blog.Name);
	}
}
func main() {
	conn, err := grpc.NewClient(addrs, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer conn.Close()
	c := pb.NewBlogServiceClient(conn)
	// id, err := doCreate(c)
	// if err == nil {

	// 	doRead(c, id)
	// }
	doList(c);
}
