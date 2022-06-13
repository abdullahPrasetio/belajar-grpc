package main

import (
	userpb "belajar-grpc/models/user"
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func serviceUser() userpb.UsersClient {
	port := ":9000"
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return userpb.NewUsersClient(conn)
}
func main() {
	user := serviceUser()
	// ctx, cancel:=context.WithCancel(context.Background())
	// defer cancel()
	// newCtx:=context.WithValue(ctx,"x-user-id","Hello")
	// fmt.Println(newCtx)]
	md := metadata.Pairs("x-user-id","Hellop ya")

	ctx := metadata.NewOutgoingContext(context.Background(), md)
	fmt.Println(ctx)
	fmt.Println(user.List(ctx,new(empty.Empty)))
}