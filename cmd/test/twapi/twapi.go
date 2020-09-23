package main

import (
  "bufio"
  "context"
  "fmt"
  "github.com/davecgh/go-spew/spew"
  "goldcrest/twitter1"
  "io/ioutil"
  "net/http"
  "os"
  "strings"
  "time"
)

func main() {
  v1()
}

func v1() {
  reader := bufio.NewReader(os.Stdin)

  consumerKey, secretKey, token, tokenSecret :=
    readLn(reader, "consumer key"),
    readLn(reader, "secret key"),
    readLn(reader, "access token"),
    readLn(reader, "token secret")

  twitter := twitter1.NewTwitter(twitter1.TwitterConfig{ClientTimeoutSeconds: 5})

  tweet, err := twitter.GetTweet(
    context.Background(),
    twitter1.Auth{Key: secretKey, Token: tokenSecret},
    twitter1.Auth{Key: consumerKey, Token: token},
    1305385801723916289,
    twitter1.DefaultTweetParams(),
  )

  if err != nil {
    panic(err)
  }

  spew.Dump(tweet)
}

func v2() {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("bearer token :> ")
  bearerToken, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  bearerToken = strings.TrimSpace(bearerToken)

  client := http.Client{
    Timeout: time.Second * 5,
  }

  req, err := http.NewRequest(
    "GET",
    //"https://api.twitter.com/2/tweets?ids=1261326399320715264,1278347468690915330",
    "https://api.twitter.com/2/tweets?ids=1304387795507650561&tweet.fields=created_at,conversation_id,attachments,entities&expansions=author_id,attachments.media_keys",
    nil,
  )
  if err != nil {
    panic(err)
  }
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  for key, value := range resp.Header {
    fmt.Println(fmt.Sprintf("%s: %s", key, value))
  }

  respBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    panic(err)
  }

  fmt.Println(string(respBody))
}

func readLn(reader *bufio.Reader, prompt string) string {
  fmt.Print(fmt.Sprintf("%s :> ", prompt))
  str, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  return strings.TrimSpace(str)
}
