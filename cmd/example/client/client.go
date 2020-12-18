package main

import (
  "context"
  "encoding/json"
  "github.com/davecgh/go-spew/spew"
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/oauth"
  "google.golang.org/grpc"
  "io/ioutil"
  "time"
)

func main() {
  authData, err := ioutil.ReadFile("conf/test-auth.json")
  if err != nil {
    panic(err)
  }
  var auth oauth.AuthPair
  if err := json.Unmarshal(authData, &auth); err != nil {
    panic(err)
  }

  authMsg := pb.Authentication{
    ConsumerKey: auth.Public.Key,
    AccessToken: auth.Public.Token,
    SecretKey:   auth.Secret.Key,
    SecretToken: auth.Secret.Token,
  }

  conn, err := grpc.Dial("localhost:7400", grpc.WithBlock(), grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := pb.NewTwitter1Client(conn)

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  tweet, err := client.GetTweet(ctx, &pb.TweetRequest{
    Auth: &authMsg,
    Id:   1339937761126658049,
  })
  if err != nil {
    panic(err)
  }

  spew.Dump(tweet)
}
