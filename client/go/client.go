package au

import (
  "context"
  "errors"
  pb "github.com/pantonshire/goldcrest/protocol"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "strconv"
  "time"
)

type Client struct {
  twitter pb.TwitterClient
  auth    authentication
  timeout time.Duration
  retry   retryPolicy
  twopts  TweetOptions
}

func NewClient(conn *grpc.ClientConn) Client {
  return Client{
    twitter: pb.NewTwitterClient(conn),
    twopts:  NewTweetOptions(),
  }
}

func (client Client) WithAuth(consumerKey, accessToken, secretKey, secretToken string) Client {
  client.auth = authentication{
    consumerKey: consumerKey,
    accessToken: accessToken,
    secretKey:   secretKey,
    secretToken: secretToken,
  }
  return client
}

func (client Client) WithTimeout(timeout time.Duration) Client {
  client.timeout = timeout
  return client
}

func (client Client) WithRetryLimit(limit uint) Client {
  client.retry = limitedRetryPolicy{
    limit: limit,
  }
  return client
}

func (client Client) WithRetryTimeout(timeout time.Duration) Client {
  client.retry = timeoutRetryPolicy{
    timeout: timeout,
  }
  return client
}

func (client Client) WithTweetOptions(twopts TweetOptions) Client {
  client.twopts = twopts
  return client
}

func (client Client) GetTweetOptions() TweetOptions {
  return client.twopts
}

type retryPolicy interface {
  newRetryer() retryer
}

type limitedRetryPolicy struct {
  limit uint
}

func (rp limitedRetryPolicy) newRetryer() retryer {
  return &limitedRetryer{
    limit: rp.limit,
  }
}

type timeoutRetryPolicy struct {
  timeout time.Duration
}

func (rp timeoutRetryPolicy) newRetryer() retryer {
  return timeoutRetryer{
    deadline: time.Now().Add(rp.timeout),
  }
}

type retryer interface {
  shouldRetry(time.Time) bool
}

type limitedRetryer struct {
  limit uint
}

func (r *limitedRetryer) shouldRetry(_ time.Time) bool {
  if r.limit == 0 {
    return false
  }
  r.limit--
  return true
}

type timeoutRetryer struct {
  deadline time.Time
}

func (r timeoutRetryer) shouldRetry(t time.Time) bool {
  return t.Before(r.deadline)
}

func (client Client) newContext() (context.Context, context.CancelFunc) {
  if client.timeout > 0 {
    return context.WithTimeout(context.Background(), client.timeout)
  }
  return context.Background(), nil
}

func (client Client) request(reqFunc func(ctx context.Context) (metadata.MD, *pb.Error, error)) error {
  rp := client.retry.newRetryer()
  for {
    meta, errMsg, err := func() (metadata.MD, *pb.Error, error) {
      ctx, cancel := client.newContext()
      if cancel != nil {
        defer cancel()
      }
      return reqFunc(ctx)
    }()
    if err != nil {
      return err
    }
    if errMsg != nil {
      if errMsg.Code == pb.Error_RATE_LIMIT {
        if meta != nil {
          if retryStrs := meta.Get("retry"); len(retryStrs) > 0 {
            retryUnix, err := strconv.ParseInt(retryStrs[0], 10, 64)
            if err != nil {
              return err
            }
            retryTime := time.Unix(retryUnix, 0)
            if rp == nil || rp.shouldRetry(retryTime) {
              time.Sleep(retryTime.Sub(time.Now())) //TODO: use a context timeout here
              continue
            }
            return RateLimitError{resets: retryTime}
          }
        }
        return AmbiguousRateLimitError{}
      } else {
        return errors.New(errMsg.Message) //TODO: wrap in custom error type
      }
    }
    return nil
  }
}

func (client Client) tweetRequest(id uint64, grpcFunc func(context.Context, *pb.TweetRequest, ...grpc.CallOption) (*pb.TweetResponse, error)) (Tweet, error) {
  var msg *pb.Tweet
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    resp, err := grpcFunc(ctx, &pb.TweetRequest{
      Auth:   client.auth.ser(),
      Id:     id,
      Twopts: client.twopts.ser(),
    }, grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.TweetResponse_Tweet); ok {
      msg = success.Tweet
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.TweetResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return Tweet{}, err
  }
  return desTweet(msg), nil
}

func (client Client) GetTweet(id uint64) (Tweet, error) {
  return client.tweetRequest(id, client.twitter.GetTweet)
}

func (client Client) GetTweets(ids ...uint64) ([]Tweet, error) {
  if len(ids) == 0 {
    return nil, nil
  }
  var msg *pb.Tweets
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    resp, err := client.twitter.GetTweets(ctx, &pb.TweetsRequest{
      Auth:   client.auth.ser(),
      Ids:    ids,
      Twopts: client.twopts.ser(),
    }, grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.TweetsResponse_Tweets); ok {
      msg = success.Tweets
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.TweetsResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return nil, err
  }
  return desTimeline(msg), nil
}

func (client Client) LikeTweet(id uint64) (Tweet, error) {
  return client.tweetRequest(id, client.twitter.LikeTweet)
}

func (client Client) UnlikeTweet(id uint64) (Tweet, error) {
  return client.tweetRequest(id, client.twitter.UnlikeTweet)
}

func (client Client) RetweetTweet(id uint64) (Tweet, error) {
  return client.tweetRequest(id, client.twitter.RetweetTweet)
}

func (client Client) UnretweetTweet(id uint64) (Tweet, error) {
  return client.tweetRequest(id, client.twitter.UnretweetTweet)
}

func (client Client) DeleteTweet(id uint64) (Tweet, error) {
  return client.tweetRequest(id, client.twitter.DeleteTweet)
}

func (client Client) HomeTimeline(tlopts TimelineOptions, replies bool) ([]Tweet, error) {
  var msg *pb.Tweets
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    resp, err := client.twitter.GetHomeTimeline(ctx, &pb.HomeTimelineRequest{
      Auth:            client.auth.ser(),
      TimelineOptions: tlopts.ser(client.twopts),
      IncludeReplies:  replies,
    }, grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.TweetsResponse_Tweets); ok {
      msg = success.Tweets
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.TweetsResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return nil, err
  }
  return desTimeline(msg), nil
}

func (client Client) MentionTimeline(tlopts TimelineOptions) ([]Tweet, error) {
  var msg *pb.Tweets
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    resp, err := client.twitter.GetMentionTimeline(ctx, &pb.MentionTimelineRequest{
      Auth:            client.auth.ser(),
      TimelineOptions: tlopts.ser(client.twopts),
    }, grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.TweetsResponse_Tweets); ok {
      msg = success.Tweets
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.TweetsResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return nil, err
  }
  return desTimeline(msg), nil
}

func (client Client) UserTimeline(user UserIdentifier, tlopts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  var msg *pb.Tweets
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    req := &pb.UserTimelineRequest{
      Auth:            client.auth.ser(),
      TimelineOptions: tlopts.ser(client.twopts),
      IncludeReplies:  replies,
      IncludeRetweets: retweets,
    }
    user.serIntoUserTimelineRequest(req)
    resp, err := client.twitter.GetUserTimeline(ctx, req, grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.TweetsResponse_Tweets); ok {
      msg = success.Tweets
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.TweetsResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return nil, err
  }
  return desTimeline(msg), nil
}

func (client Client) PublishTweet(com TweetComposer) (Tweet, error) {
  var msg *pb.Tweet
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    resp, err := client.twitter.PublishTweet(ctx, com.ser(client.auth, client.twopts), grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.TweetResponse_Tweet); ok {
      msg = success.Tweet
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.TweetResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return Tweet{}, err
  }
  return desTweet(msg), nil
}

func (client Client) UpdateProfile(pu ProfileUpdater, includeEntities, includeStatuses bool) (User, error) {
  var msg *pb.User
  err := client.request(func(ctx context.Context) (metadata.MD, *pb.Error, error) {
    var header metadata.MD
    resp, err := client.twitter.UpdateProfile(ctx, pu.ser(client.auth, includeEntities, includeStatuses), grpc.Header(&header))
    if err != nil {
      return nil, nil, err
    }
    if success, ok := resp.Response.(*pb.UserResponse_User); ok {
      msg = success.User
      return header, nil, nil
    } else if failure, ok := resp.Response.(*pb.UserResponse_Error); ok {
      return header, failure.Error, nil
    } else {
      return header, nil, errors.New("invalid response")
    }
  })
  if err != nil {
    return User{}, err
  }
  return desUser(msg), nil
}
