package twitter1

import (
  "context"
  "goldcrest/rpc"
  "google.golang.org/grpc"
  "path"
)

type twitterServer struct {
  twitter *Twitter
}

func (t *Twitter) Server(server *grpc.Server) error {
  ts := &twitterServer{twitter: t}
  rpc.RegisterTwitter1Server(server, ts)
  return nil
}

func (s *twitterServer) GetTweet(ctx context.Context, req *rpc.TweetRequest) (*rpc.Tweet, error) {
  secret, auth := decodeAuthPair(req.Auth)
  opts := decodeTweetOptions(req.Options)
  mod, err := s.twitter.GetTweet(secret, auth, req.Id, opts)
  if err != nil {
    return nil, err
  }
  return tweetModelToMessage(mod)
}

func (s *twitterServer) GetTweets(req *rpc.TweetsRequest, srv rpc.Twitter1_GetTweetsServer) error {
  return nil
}

func (s *twitterServer) GetRaw(ctx context.Context, rr *rpc.RawAPIRequest) (*rpc.RawAPIResult, error) {
  secret, auth := decodeAuthPair(rr.Auth)
  or := OAuthRequest{
    Method:   rr.Method,
    Protocol: rr.Protocol,
    Domain:   domain,
    Path:     path.Join(rr.Version, rr.Path),
    Query:    rr.QueryParams,
    Body:     rr.BodyParams,
  }
  req, err := or.MakeRequest(secret, auth)
  if err != nil {
    return nil, err
  }
  status, headers, body, err := s.twitter.requestRaw(req)
  if err != nil {
    return nil, err
  }
  return &rpc.RawAPIResult{Status: uint32(status), Headers: headers, Body: body}, nil
}
