package main

import (
  "context"
  "encoding/json"
  "github.com/davecgh/go-spew/spew"
  "github.com/pantonshire/goldcrest/client/go/au"
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

  ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
  defer cancel()
  conn, err := grpc.DialContext(ctx, "localhost:7400", grpc.WithBlock(), grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := au.NewClient(conn).
    WithAuth(auth.Public.Key, auth.Public.Token, auth.Secret.Key, auth.Secret.Token).
    WithTimeout(time.Second * 5).
    WithRetryLimit(5)

  resp, err := client.GetTweet(1224097454477512704)
  //resp, err := client.RetweetTweet(1309868095306162176)
  //resp, err := client.UserTimeline(au.UserHandle("smolrobots"), au.NewTimelineOptions(5), true, true)
  //com := au.NewTweetComposer("Reply test :D").
  //  ReplyTo(1341792045288091648)
  //resp, err := client.PublishTweet(com)
  //resp, err := client.WithTweetOptions(client.GetTweetOptions().WithTrimUser(true)).
  //  HomeTimeline(au.NewTimelineOptions(6), true)
  if err != nil {
    panic(err)
  }

  spew.Dump(resp)
}
