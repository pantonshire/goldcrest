package twitter1

import (
  "context"
  "goldcrest"
  pb "goldcrest/proto"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "time"
)

type Client interface {
  GetTweet(params TweetOptions, id uint64) (Tweet, error)
}

type local struct {
  secret, auth Auth
}

func Local(secret, auth Auth) Client {
  return local{secret: secret, auth: auth}
}

func (lc local) GetTweet(params TweetOptions, id uint64) (Tweet, error) {
  return Tweet{}, nil
}

//TODO: server health checks
type remote struct {
  secret, auth Auth
  address      string
  client       pb.Twitter1Client
  callTimeout  time.Duration
}

func Remote(conn *grpc.ClientConn, secret, auth Auth, timeout time.Duration) Client {
  return remote{
    secret:      secret,
    auth:        auth,
    address:     conn.Target(),
    client:      pb.NewTwitter1Client(conn),
    callTimeout: timeout,
  }
}

func (rc remote) newContext() (context.Context, context.CancelFunc) {
  if rc.callTimeout == 0 {
    return context.Background(), nil
  }
  return context.WithTimeout(context.Background(), rc.callTimeout)
}

func (rc remote) handleRequest(handler func() error) error {
  err := handler()
  if httpErr, ok := err.(*goldcrest.HttpError); ok {
    return status.Errorf(codes.Internal, "twitter error %s", httpErr.Error())
  }
  return err
}

func (rc remote) GetTweet(params TweetOptions, id uint64) (tweet Tweet, err error) {
  err = rc.handleRequest(func() error {
    ctx, cancel := rc.newContext()
    if cancel != nil {
      defer cancel()
    }
    tweetMsg, err := rc.client.GetTweet(ctx, &pb.TweetRequest{
      Auth:    encodeAuthPair(rc.secret, rc.auth),
      Id:      id,
      Options: encodeTweetOptions(params),
    })
    if err != nil {
      return err
    }
    tweet = decodeTweet(tweetMsg)
    return nil
  })
  if err != nil {
    return Tweet{}, err
  }
  return tweet, nil
}

func (rc remote) GetRaw(method, protocol, version, path string, queryParams, bodyParams map[string]string) (headers map[string]string, status uint, body []byte, err error) {
  err = rc.handleRequest(func() error {
    ctx, cancel := rc.newContext()
    if cancel != nil {
      defer cancel()
    }
    resp, err := rc.client.GetRaw(ctx, &pb.RawAPIRequest{
      Auth:        encodeAuthPair(rc.secret, rc.auth),
      Method:      method,
      Protocol:    protocol,
      Version:     version,
      Path:        path,
      QueryParams: queryParams,
      BodyParams:  bodyParams,
    })
    if err != nil {
      return err
    }
    headers = resp.Headers
    status = uint(resp.Status)
    body = resp.Body
    return nil
  })
  if err != nil {
    return nil, 0, nil, err
  }
  return headers, status, body, nil
}
