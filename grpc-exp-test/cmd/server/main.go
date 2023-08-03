package main

import (
	"context"
	"fmt"
	ms "grpc-exp-test/internal/protos/message_service"
	"net"

	"google.golang.org/grpc"
)

type MessageServer struct {
}

func (MessageServer) SayIt(ctx context.Context, r *ms.Request) (*ms.Response, error) {
	fmt.Println("Request Text:", r.Text)
	fmt.Println("Request SubText:", r.Subtext)
	response := &ms.Response{
		Text:    r.Text,
		Subtext: "Got it!",
	}

	return response, nil
}

var messageServer MessageServer
var port = ":8080"

func main() {
	server := grpc.NewServer()
	ms.RegisterMessageServiceServer(server, messageServer)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Serving requests...")
	server.Serve(listen)
}
