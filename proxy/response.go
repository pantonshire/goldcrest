package proxy

import (
  pb "github.com/pantonshire/goldcrest/protocol"
  "google.golang.org/grpc/metadata"
)

func serError(err error) (*pb.Error, metadata.MD) {
  if proxyError, ok := err.(proxyError); ok {
    return proxyError.ser()
  }
  return nil, nil
}

func generateTweetResponse(generator func() (*pb.Tweet, metadata.MD, error)) (*pb.TweetResponse, metadata.MD, error) {
  tweet, meta, err := generator()
  if err != nil {
    if errMsg, errMeta := serError(err); errMsg != nil {
      return &pb.TweetResponse{Response: &pb.TweetResponse_Error{Error: errMsg}}, metadata.Join(meta, errMeta), nil
    }
    return nil, nil, err
  }
  return &pb.TweetResponse{Response: &pb.TweetResponse_Tweet{Tweet: tweet}}, meta, nil
}

func generateTimelineResponse(generator func() (*pb.Timeline, metadata.MD, error)) (*pb.TimelineResponse, metadata.MD, error) {
  timeline, meta, err := generator()
  if err != nil {
    if errMsg, errMeta := serError(err); errMsg != nil {
      return &pb.TimelineResponse{Response: &pb.TimelineResponse_Error{Error: errMsg}}, metadata.Join(meta, errMeta), nil
    }
    return nil, nil, err
  }
  return &pb.TimelineResponse{Response: &pb.TimelineResponse_Timeline{Timeline: timeline}}, meta, nil
}

func generateUserResponse(generator func() (*pb.User, metadata.MD, error)) (*pb.UserResponse, metadata.MD, error) {
  user, meta, err := generator()
  if err != nil {
    if errMsg, errMeta := serError(err); errMsg != nil {
      return &pb.UserResponse{Response: &pb.UserResponse_Error{Error: errMsg}}, metadata.Join(meta, errMeta), nil
    }
    return nil, nil, err
  }
  return &pb.UserResponse{Response: &pb.UserResponse_User{User: user}}, meta, nil
}

