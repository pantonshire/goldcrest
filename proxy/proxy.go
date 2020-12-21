package proxy

import (
  "context"
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/model"
  "github.com/pantonshire/goldcrest/proxy/oauth"
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "strconv"
)

var log = logrus.New()

type proxy struct {
  tc twitterClient
}

func InitServer(server *grpc.Server) {
  log.SetLevel(logrus.DebugLevel)

  p := proxy{
    tc: newTwitterClient(),
  }

  pb.RegisterTwitter1Server(server, &p)
}

func (p proxy) GetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(showTweetEndpoint, auth, query, nil, &tweet); err != nil {
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
  auth, query := reserTweetsRequest(req)
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    var tweets model.Timeline
    if err := p.tc.standardRequest(showTweetsEndpoint, auth, query, nil, &tweets); err != nil {
      return nil, nil, err
    }
    return tweets, nil, nil
  })
  if err != nil {
    return nil, err
  }
  if err := sendHeader(ctx, meta); err != nil {
    return nil, err
  }
  return resp, nil
}

func (p proxy) LikeTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(likeEndpoint, auth, query, nil, &tweet); err != nil {
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

func (p proxy) UnlikeTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(unlikeEndpoint, auth, query, nil, &tweet); err != nil {
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

func (p proxy) RetweetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(retweetEndpoint, auth, query, nil, &tweet); err != nil {
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

func (p proxy) UnretweetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(unretweetEndpoint, auth, query, nil, &tweet); err != nil {
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

func (p proxy) PublishTweet(ctx context.Context, req *pb.PublishTweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserPublishTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(publishTweetEndpoint, auth, query, nil, &tweet); err != nil {
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

func (p proxy) DeleteTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth, query := reserTweetRequest(req)
  resp, meta, err := generateTweetResponse(func() (model.Tweet, metadata.MD, error) {
    var tweet model.Tweet
    if err := p.tc.standardRequest(destroyTweetEndpoint, auth, query, nil, &tweet); err != nil {
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

func (p proxy) GetHomeTimeline(ctx context.Context, req *pb.HomeTimelineRequest) (*pb.TweetsResponse, error) {
  auth := desAuth(req.GetAuth())
  query := desTimelineOptions(req.GetTimelineOptions()).ser()
  query.Set("exclude_replies", strconv.FormatBool(!req.GetIncludeReplies()))
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    var tweets model.Timeline
    if err := p.tc.standardRequest(homeTimelineEndpoint, auth, query, nil, &tweets); err != nil {
      return nil, nil, err
    }
    return tweets, nil, nil
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
  auth := desAuth(req.GetAuth())
  query := desTimelineOptions(req.GetTimelineOptions()).ser()
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    var tweets model.Timeline
    if err := p.tc.standardRequest(mentionTimelineEndpoint, auth, query, nil, &tweets); err != nil {
      return nil, nil, err
    }
    return tweets, nil, nil
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
  auth := desAuth(req.GetAuth())
  query := desTimelineOptions(req.GetTimelineOptions()).ser()
  if id, ok := req.GetUser().(*pb.UserTimelineRequest_UserId); ok {
    query.Set("user_id", strconv.FormatUint(id.UserId, 10))
  } else if handle, ok := req.GetUser().(*pb.UserTimelineRequest_UserHandle); ok {
    query.Set("screen_name", handle.UserHandle)
  }
  query.Set("exclude_replies", strconv.FormatBool(!req.GetIncludeReplies()))
  query.Set("include_rts", strconv.FormatBool(req.GetIncludeRetweets()))
  resp, meta, err := generateTweetsResponse(func() (model.Timeline, metadata.MD, error) {
    var tweets model.Timeline
    if err := p.tc.standardRequest(userTimelineEndpoint, auth, query, nil, &tweets); err != nil {
      return nil, nil, err
    }
    return tweets, nil, nil
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
  auth := desAuth(req.GetAuth())
  query := oauth.NewParams()
  query.Set("include_entities", strconv.FormatBool(req.GetIncludeEntities()))
  query.Set("skip_status", strconv.FormatBool(!req.GetIncludeStatuses()))
  if name, ok := req.GetUpdateName().(*pb.UpdateProfileRequest_Name); ok {
    query.Set("name", name.Name)
  }
  if url, ok := req.GetUpdateUrl().(*pb.UpdateProfileRequest_Url); ok {
    query.Set("url", url.Url)
  }
  if location, ok := req.GetUpdateLocation().(*pb.UpdateProfileRequest_Location); ok {
    query.Set("location", location.Location)
  }
  if bio, ok := req.GetUpdateBio().(*pb.UpdateProfileRequest_Bio); ok {
    query.Set("description", bio.Bio)
  }
  if col, ok := req.GetUpdateProfileLinkColor().(*pb.UpdateProfileRequest_ProfileLinkColor); ok {
    query.Set("profile_link_color", col.ProfileLinkColor)
  }
  resp, meta, err := generateUserResponse(func() (model.User, metadata.MD, error) {
    var user model.User
    if err := p.tc.standardRequest(updateProfileEndpoint, auth, query, nil, &user); err != nil {
      return model.User{}, nil, err
    }
    return user, nil, nil
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
