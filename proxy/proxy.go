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
  "time"
)

var log *logrus.Logger

type Proxy struct {
  tc twitterClient
}

func NewProxy(logger *logrus.Logger, twitterTimeout time.Duration, twitterProtocol, twitterURL string, assumeNextLimit bool) *Proxy {
  log = logger
  return &Proxy{
    tc: newTwitterClient(twitterTimeout, twitterProtocol, twitterURL, assumeNextLimit),
  }
}

func (p Proxy) GetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) GetTweets(ctx context.Context, req *pb.TweetsRequest) (*pb.TweetsResponse, error) {
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

func (p Proxy) LikeTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) UnlikeTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) RetweetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) UnretweetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) PublishTweet(ctx context.Context, req *pb.PublishTweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) DeleteTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
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

func (p Proxy) GetHomeTimeline(ctx context.Context, req *pb.HomeTimelineRequest) (*pb.TweetsResponse, error) {
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

func (p Proxy) GetMentionTimeline(ctx context.Context, req *pb.MentionTimelineRequest) (*pb.TweetsResponse, error) {
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

func (p Proxy) GetUserTimeline(ctx context.Context, req *pb.UserTimelineRequest) (*pb.TweetsResponse, error) {
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

func (p Proxy) SearchTweets(ctx context.Context, req *pb.SearchRequest) (*pb.TweetsResponse, error) {
  auth := desAuth(req.GetAuth())

  query := desTimelineOptions(req.GetTimelineOptions()).ser()
  query.Set("q", req.GetQuery())
  if geo := req.GetGeocode(); geo != nil {
    query.Set("geocode", geo.Val)
  }
  if lang := req.GetLang(); lang != nil {
    query.Set("lang", lang.Val)
  }
  if locale := req.GetLocale(); locale != nil {
    query.Set("locale", locale.Val)
  }
  query.Set("result_type", reserSearchResultType(req.GetResultType()))
  if untilUnix := req.GetUntilTimestamp(); untilUnix != nil {
    until := time.Unix(int64(untilUnix.Val), 0)
    untilStr := until.Format("2006-01-02")
    query.Set("until", untilStr)
  }

  resp, meta, err := generateSearchResultResponse(func() (model.SearchResult, metadata.MD, error) {
    var tweets model.SearchResult
    if err := p.tc.standardRequest(searchEndpoint, auth, query, nil, &tweets); err != nil {
      return model.SearchResult{}, nil, err
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

func (p Proxy) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
  auth := desAuth(req.GetAuth())
  query := oauth.NewParams()
  query.Set("include_entities", strconv.FormatBool(req.GetIncludeEntities()))
  query.Set("skip_status", strconv.FormatBool(!req.GetIncludeStatuses()))
  if req.Name != nil {
    query.Set("name", req.Name.Val)
  }
  if req.Url != nil {
    query.Set("url", req.Url.Val)
  }
  if req.Location != nil {
    query.Set("location", req.Location.Val)
  }
  if req.Bio != nil {
    query.Set("description", req.Bio.Val)
  }
  if req.LinkColor != nil {
    query.Set("profile_link_color", req.LinkColor.Val)
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

func (p Proxy) GetRaw(ctx context.Context, req *pb.RawAPIRequest) (*pb.RawAPIResult, error) {
  //TODO
  panic("implement me")
}

func sendHeader(ctx context.Context, meta metadata.MD) error {
  if meta != nil {
    return grpc.SendHeader(ctx, meta)
  }
  return nil
}
