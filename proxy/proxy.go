package proxy

import (
  "context"
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/model"
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
)

var log = logrus.New()

type proxy struct {
  tc twitterClient
}

func InitServer(server *grpc.Server) {
  log.SetLevel(logrus.DebugLevel)
  pb.RegisterTwitter1Server(server, &proxy{tc: newTwitterClient()})
}

func (p proxy) GetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, id, opts := desTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    //TODO: potentially use a different context with a custom timeout?
    tweet, err := p.tc.getTweet(ctx, auth, id, opts)
    if err != nil {
      return model.Tweet{}, nil, err
    }
    return tweet, nil, nil
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
}

func (p proxy) GetTweets(ctx context.Context, req *pb.TweetsRequest) (*pb.TweetsResponse, error) {
  panic("implement me")
}

func (p proxy) LikeTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  panic("implement me")
}

func (p proxy) UnlikeTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  panic("implement me")
}

func (p proxy) RetweetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  panic("implement me")
}

func (p proxy) UnretweetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  panic("implement me")
}

func (p proxy) DeleteTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  panic("implement me")
}

func (p proxy) GetHomeTimeline(ctx context.Context, req *pb.HomeTimelineRequest) (*pb.TweetsResponse, error) {
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
}

func (p proxy) GetMentionTimeline(ctx context.Context, req *pb.MentionTimelineRequest) (*pb.TweetsResponse, error) {
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
}

func (p proxy) GetUserTimeline(ctx context.Context, req *pb.UserTimelineRequest) (*pb.TweetsResponse, error) {
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
}

func (p proxy) PublishTweet(ctx context.Context, req *pb.PublishTweetRequest) (*pb.TweetResponse, error) {
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
}

func (p proxy) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
  resp, meta, err := generateUserResponse(func() (model.User, metadata.MD, error) {
    //TODO
    panic("implement me")
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
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
