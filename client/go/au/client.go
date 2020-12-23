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
              time.Sleep(time.Now().Sub(retryTime)) //TODO: use a context timeout here
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
