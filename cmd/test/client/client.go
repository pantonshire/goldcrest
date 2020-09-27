package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "github.com/davecgh/go-spew/spew"
  "goldcrest/twitter1"
  "google.golang.org/grpc"
  "io/ioutil"
  "strings"
  "time"
)

func main() {
  authData, err := ioutil.ReadFile("conf/test-auth.json")
  if err != nil {
    panic(err)
  }
  var auth twitter1.AuthPair
  if err := json.Unmarshal(authData, &auth); err != nil {
    panic(err)
  }

  conn, err := grpc.Dial("localhost:7400", grpc.WithBlock(), grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := twitter1.Remote(conn, auth.Secret, auth.Public, time.Second*5)

  //tweet, err := client.GetTweet(twitter1.DefaultTweetOptions(), 1305748179338629120)

  //var replyID uint64 = 1309875852180705282
  //tweet, err := client.UpdateStatus("@SmolbotbotT Reply test", twitter1.StatusUpdateOptions{ReplyID: &replyID}, false)

  //timeline, err := client.GetUserHandleTimeline(
  //  twitter1.TweetOptions{
  //    TrimUser:          true,
  //    IncludeMyRetweet:  true,
  //    IncludeEntities:   true,
  //    IncludeExtAltText: true,
  //    IncludeCardURI:    true,
  //    Mode:              twitter1.ExtendedMode,
  //  },
  //  "smolrobots",
  //  twitter1.TimelineOptions{Count: 3},
  //  true,
  //  true,
  //)

  url := "github.com/pantonshire/smolbotbot"
  location := "Test"

  user, err := client.UpdateProfile(twitter1.ProfileUpdateOptions{
    Url:      &url,
    Location: &location,
  }, true, true)

  if err != nil {
    panic(err)
  }

  //spew.Dump(tweet)
  //
  //fmt.Println(tweet.TextOnly())
  //if tweet.Quoted != nil {
  //  fmt.Println(tweet.Quoted.TextOnly())
  //}

  //spew.Dump(timeline)

  spew.Dump(user)
}

func readLn(reader *bufio.Reader, prompt string) string {
  fmt.Print(fmt.Sprintf("%s :> ", prompt))
  str, err := reader.ReadString('\n')
  if err != nil {
    panic(err)
  }
  return strings.TrimSpace(str)
}
