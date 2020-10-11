package main

import (
  "encoding/json"
  "fmt"
  "github.com/davecgh/go-spew/spew"
  "goldcrest/twitter1"
  "io/ioutil"
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

  client := twitter1.Local(auth, twitter1.TwitterConfig{ClientTimeoutSeconds: 5}, time.Second*5)

  tweet, err := client.GetTweet(twitter1.DefaultTweetOptions(), 1305748179338629120)

  if err != nil {
    panic(err)
  }

  spew.Dump(tweet)
  fmt.Println(tweet.TextOnly())
  if tweet.Quoted != nil {
   fmt.Println(tweet.Quoted.TextOnly())
  }
}
