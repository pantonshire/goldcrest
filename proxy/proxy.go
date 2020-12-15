package proxy

import (
  "context"
  pb "github.com/pantonshire/goldcrest/protocol"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
)

type proxy struct {
  client twitterClient
}

func (p proxy) GetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  rsp, meta, err := generateTweetResponse(func() (*pb.Tweet, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return rsp, nil
}

func (p proxy) GetHomeTimeline(ctx context.Context, req *pb.HomeTimelineRequest) (*pb.TimelineResponse, error) {
  rsp, meta, err := generateTimelineResponse(func() (*pb.Timeline, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return rsp, nil
}

func (p proxy) GetMentionTimeline(ctx context.Context, req *pb.MentionTimelineRequest) (*pb.TimelineResponse, error) {
  rsp, meta, err := generateTimelineResponse(func() (*pb.Timeline, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return rsp, nil
}

func (p proxy) GetUserTimeline(ctx context.Context, req *pb.UserTimelineRequest) (*pb.TimelineResponse, error) {
  rsp, meta, err := generateTimelineResponse(func() (*pb.Timeline, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return rsp, nil
}

func (p proxy) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.TweetResponse, error) {
  rsp, meta, err := generateTweetResponse(func() (*pb.Tweet, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return rsp, nil
}

func (p proxy) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
  rsp, meta, err := generateUserResponse(func() (*pb.User, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return rsp, nil
}

func (p proxy) GetRaw(ctx context.Context, req *pb.RawAPIRequest) (*pb.RawAPIResult, error) {
  //TODO
  panic("implement me")
}

func sendHeader(ctx context.Context, meta metadata.MD) error {
  if meta != nil {
    return grpc.SendHeader(ctx, meta)
  }
  return nil
}
