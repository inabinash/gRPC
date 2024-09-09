package main

// import (
// 	pb "github.com/inabinash/grpc/blog/proto"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type BlogItem struct {
// 	Id      primitive.ObjectID `bson:"_id,omitempty"`
// 	Name    string             `bson:"name"`
// 	Author  string             `bson:"author"`
// 	Content string             `bson:"content"`
// }

// func DocumentToBlog(data *BlogItem) *pb.Blog {
// 	return &pb.Blog{
// 		Id:      data.Id.Hex(),
// 		Author:  data.Author,
// 		Name:    data.Name,
// 		Content: data.Content,
// 	}
// }
