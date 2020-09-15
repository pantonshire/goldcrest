package main

import (
  "fmt"
  "goldcrest/twitter1"
  "google.golang.org/grpc"
  "net"
)

func main() {
  listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 7400))
  if err != nil {
    panic(err)
  }
  grpcServer := grpc.NewServer()

  twitter := twitter1.NewTwitter(twitter1.TwitterConfig{ClientTimeoutSeconds: 5})

  if err := twitter.Server(grpcServer); err != nil {
    panic(err)
  }

  if err := grpcServer.Serve(listener); err != nil {
    panic(err)
  }
}
