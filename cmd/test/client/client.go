package main

import (
  "bufio"
  "fmt"
  "github.com/davecgh/go-spew/spew"
  "goldcrest/twitter1"
  "google.golang.org/grpc"
  "os"
  "strings"
)

func main() {
  reader := bufio.NewReader(os.Stdin)

  consumerKey, secretKey, token, tokenSecret :=
    readLn(reader, "consumer key"),
    readLn(reader, "secret key"),
    readLn(reader, "access token"),
    readLn(reader, "token secret")

  conn, err := grpc.Dial("localhost:7400", grpc.WithBlock(), grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := twitter1.Remote(conn,
    twitter1.Auth{Key: secretKey, Token: tokenSecret},
    twitter1.Auth{Key: consumerKey, Token: token},
  )

  tweet, err := client.GetTweet(twitter1.DefaultTweetParams(), 1305748179338629120)

  if err != nil {
    panic(err)
  }

  spew.Dump(tweet)

  fmt.Println(tweet.TextOnly())
  if tweet.Quoted != nil {
    fmt.Println(tweet.Quoted.TextOnly())
  }
}

func readLn(reader *bufio.Reader, prompt string) string {
  fmt.Print(fmt.Sprintf("%s :> ", prompt))
  str, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  return strings.TrimSpace(str)
}
