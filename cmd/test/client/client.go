package main

import (
  "bufio"
  "context"
  "fmt"
  "github.com/davecgh/go-spew/spew"
  "goldcrest/rpc"
  "google.golang.org/grpc"
  "os"
  "strings"
)

func main() {
  conn, err := grpc.Dial("localhost:7400", grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := rpc.NewTwitter1Client(conn)

  reader := bufio.NewReader(os.Stdin)

  consumerKey, secretKey, token, tokenSecret :=
    readLn(reader, "consumer key"),
    readLn(reader, "secret key"),
    readLn(reader, "access token"),
    readLn(reader, "token secret")

  tweet, err := client.GetTweet(context.Background(), &rpc.TweetRequest{
    Auth: &rpc.Authentication{
      ConsumerKey: consumerKey,
      AccessToken: token,
      SecretKey:   secretKey,
      SecretToken: tokenSecret,
    },
    Id: 1305748179338629120,
    Options: &rpc.TweetOptions{
      TrimUser:          false,
      IncludeMyRetweet:  true,
      IncludeEntities:   true,
      IncludeExtAltText: true,
      IncludeCardUri:    true,
      Mode:              rpc.TweetOptions_EXTENDED,
    },
  })

  if err != nil {
    panic(err)
  }

  spew.Dump(tweet)
}

func readLn(reader *bufio.Reader, prompt string) string {
  fmt.Print(fmt.Sprintf("%s :> ", prompt))
  str, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  return strings.TrimSpace(str)
}
