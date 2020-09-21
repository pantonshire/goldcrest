package twitter1

import (
  "context"
  "goldcrest"
  pb "goldcrest/proto"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

type Client interface {
  GetTweet(params TweetParams, id uint64) (Tweet, error)
}

type local struct {
  secret, auth Auth
}

func Local(secret, auth Auth) Client {
  return local{secret: secret, auth: auth}
}

func (lc local) GetTweet(params TweetParams, id uint64) (Tweet, error) {
  return Tweet{}, nil
}

//TODO: server health checks
type remote struct {
  secret, auth Auth
  address      string
  client       pb.Twitter1Client
}

func Remote(conn *grpc.ClientConn, secret, auth Auth) Client {
  return remote{
    secret:  secret,
    auth:    auth,
    address: conn.Target(),
    client:  pb.NewTwitter1Client(conn),
  }
}

func (rc remote) handleRequest(handler func() error) error {
  err := handler()
  if httpErr, ok := err.(*goldcrest.HttpError); ok {
    return status.Errorf(codes.Internal, "twitter error %s", httpErr.Error())
  }
  return err
}

func (rc remote) GetTweet(params TweetParams, id uint64) (Tweet, error) {
  var tweet Tweet
  err := rc.handleRequest(func() error {
    tweetMsg, err := rc.client.GetTweet(context.Background(), &pb.TweetRequest{
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
