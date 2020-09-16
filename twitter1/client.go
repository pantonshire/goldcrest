package twitter1

import (
  "context"
  "goldcrest/rpc"
  "google.golang.org/grpc"
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

//TODO: connection pool to reuse connections
//TODO: server health checks
type remote struct {
  secret, auth Auth
  address      string
}

func Remote(secret, auth Auth, address string) (client Client, closeClient func() error) {
  client = remote{
    secret:  secret,
    auth:    auth,
    address: address,
  }

  closeClient = func() error {
    return nil
  }

  return client, closeClient
}

//TODO: do not open a new connection each time
//TODO: connection options (eg. TLS)
func (rc remote) GetTweet(params TweetParams, id uint64) (Tweet, error) {
  conn, err := grpc.Dial(rc.address, grpc.WithInsecure())
  if err != nil {
    return Tweet{}, err
  }
  defer conn.Close()

  client := rpc.NewTwitter1Client(conn)

  tweetMsg, err := client.GetTweet(context.Background(), &rpc.TweetRequest{
    Auth:    encodeAuthPair(rc.secret, rc.auth),
    Id:      id,
    Options: encodeTweetOptions(params),
  })

  if err != nil {
    return Tweet{}, err
  }

  tweet := decodeTweet(tweetMsg)

  return tweet, nil
}
