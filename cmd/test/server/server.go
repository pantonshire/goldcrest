package main

import (
  "context"
  "fmt"
  "goldcrest/rpc"
  "google.golang.org/grpc"
  "net"
)

func main() {
  listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 7400))
  if err != nil {
    panic(err)
  }
  grpcServer := grpc.NewServer()
  rpc.RegisterTestServer(grpcServer, &server{})
  if err := grpcServer.Serve(listener); err != nil {
    panic(err)
  }
}

type server struct {}

func (s *server) GetTweet(ctx context.Context, twid *rpc.TweetID) (*rpc.Tweet, error) {
  return &rpc.Tweet{
    Id:   twid.Id,
    Text: "foo baa",
    Author: &rpc.User{
      Id:     123,
      Handle: "@tom",
      Name:   "Tom P",
    }}, nil
}
