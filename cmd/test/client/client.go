package main

import (
  "context"
  "fmt"
  "goldcrest/rpc"
  "google.golang.org/grpc"
)

func main() {
  conn, err := grpc.Dial("localhost:7400", grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := rpc.NewTestClient(conn)

  tweet, err := client.GetTweet(context.Background(), &rpc.TweetID{Id: 1000})
  if err != nil {
    panic(err)
  }

  fmt.Println(tweet)
}
