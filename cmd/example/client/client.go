package main

import (
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

  //authMsg := pb.Authentication{
  //  ConsumerKey: auth.Public.Key,
  //  AccessToken: auth.Public.Token,
  //  SecretKey:   auth.Secret.Key,
  //  SecretToken: auth.Secret.Token,
  //}

  conn, err := grpc.Dial("localhost:7400", grpc.WithBlock(), grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  defer conn.Close()

  client := au.NewClient(conn).
    WithAuth(auth.Public.Key, auth.Public.Token, auth.Secret.Key, auth.Secret.Token).
    WithTimeout(time.Second * 5).
    WithRetryLimit(5)

  //resp, err := client.GetTweet(1224097454477512704)
  //resp, err := client.RetweetTweet(1309868095306162176)
  resp, err := client.UserTimeline(au.UserHandle("smolrobots"), au.NewTimelineOptions(5), true, true)
  if err != nil {
    panic(err)
  }

  spew.Dump(resp)

  //client := pb.NewTwitterClient(conn)
  //
  //ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  //defer cancel()

  //resp, err := client.GetTweets(ctx, &pb.TweetsRequest{
  //  Auth: &authMsg,
  //  Ids:  []uint64{1339937761126658049, 1250097187377360896, 1339877652929245184, 1074568259385597952},
  //})

  //resp, err := client.PublishTweet(ctx, &pb.PublishTweetRequest{
  // Auth:                      &authMsg,
  // Text:                      "I think I messed up my rate limit tracker?",
  // AutoPopulateReplyMetadata: false,
  // PossiblySensitive:         false,
  //})

  //resp, err := client.RetweetTweet(ctx, &pb.TweetRequest{
  // Auth: &authMsg,
  // Id:   1340711006440480768,
  //})

  //resp, err := client.GetUserTimeline(ctx, &pb.UserTimelineRequest{
  //  Auth: &authMsg,
  //  User: &pb.UserTimelineRequest_UserHandle{UserHandle: "PantonshireDev"},
  //  TimelineOptions: &pb.TimelineOptions{
  //    Count: 15,
  //    Content: &pb.TimelineOptions_Custom{
  //      Custom: &pb.TweetOptions{
  //        TrimUser:          true,
  //        IncludeMyRetweet:  true,
  //        IncludeEntities:   true,
  //        IncludeExtAltText: true,
  //        IncludeCardUri:    true,
  //        Mode:              pb.TweetOptions_EXTENDED,
  //      },
  //    },
  //  },
  //  IncludeReplies:  true,
  //  IncludeRetweets: true,
  //})

  //resp, err := client.UpdateProfile(ctx, &pb.UpdateProfileRequest{
  //  Auth:           &authMsg,
  //  UpdateLocation: &pb.UpdateProfileRequest_Location{Location: "Testing with Goldcrest v0.1"},
  //  UpdateBio:      &pb.UpdateProfileRequest_Bio{Bio: "Test account for @smolbotbot!"},
  //  UpdateName:     &pb.UpdateProfileRequest_Name{Name: "Smolbotbot Test"},
  //})
  //
  //if err != nil {
  //  panic(err)
  //}
  //
  //s, err := json.MarshalIndent(resp, "", " ")
  //if err != nil {
  //  panic(err)
  //}
  //fmt.Println(string(s))
}
