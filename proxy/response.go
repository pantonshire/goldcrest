package proxy

import (
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/model"
  "google.golang.org/grpc/metadata"
)

func serError(err error) (*pb.Error, metadata.MD) {
  if proxyError, ok := err.(proxyError); ok {
    return proxyError.ser()
  }
  return nil, nil
}

func generateTweetResponse(generator func() (model.Tweet, metadata.MD, error)) (*pb.TweetResponse, metadata.MD, error) {
  tweet, meta, err := generator()
  if err != nil {
    if errMsg, errMeta := serError(err); errMsg != nil {
      return &pb.TweetResponse{Response: &pb.TweetResponse_Error{Error: errMsg}}, metadata.Join(meta, errMeta), nil
    }
    return nil, nil, err
  }
  return &pb.TweetResponse{Response: &pb.TweetResponse_Tweet{Tweet: serTweet(tweet)}}, meta, nil
}

func generateTweetsResponse(generator func() (model.Timeline, metadata.MD, error)) (*pb.TweetsResponse, metadata.MD, error) {
  tweets, meta, err := generator()
  if err != nil {
    if errMsg, errMeta := serError(err); errMsg != nil {
      return &pb.TweetsResponse{Response: &pb.TweetsResponse_Error{Error: errMsg}}, metadata.Join(meta, errMeta), nil
    }
    return nil, nil, err
  }
  return &pb.TweetsResponse{Response: &pb.TweetsResponse_Tweets{Tweets: serTimeline(tweets)}}, meta, nil
}

func generateUserResponse(generator func() (model.User, metadata.MD, error)) (*pb.UserResponse, metadata.MD, error) {
  user, meta, err := generator()
  if err != nil {
    if errMsg, errMeta := serError(err); errMsg != nil {
      return &pb.UserResponse{Response: &pb.UserResponse_Error{Error: errMsg}}, metadata.Join(meta, errMeta), nil
    }
    return nil, nil, err
  }
  return &pb.UserResponse{Response: &pb.UserResponse_User{User: serUser(user)}}, meta, nil
}

