package twitter1

import (
  "context"
  pb "goldcrest/proto"
  "google.golang.org/grpc"
  "path"
)

type twitterServer struct {
  twitter *Twitter
}

func (t *Twitter) Server(server *grpc.Server) error {
  ts := &twitterServer{twitter: t}
  pb.RegisterTwitter1Server(server, ts)
  return nil
}

func (s *twitterServer) GetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.Tweet, error) {
  auth := decodeAuthPair(req.Auth)
  opts := decodeTweetOptions(req.Options)
  mod, err := s.twitter.GetTweet(ctx, auth, req.Id, opts)
  if err != nil {
    return nil, err
  }
  return tweetModelToMessage(mod)
}

func (s *twitterServer) GetTweets(req *pb.TweetsRequest, srv pb.Twitter1_GetTweetsServer) error {
  return nil
}

func (s *twitterServer) GetHomeTimeline(ctx context.Context, req *pb.HomeTimelineRequest) (*pb.Timeline, error) {
  panic("implement me")
}

func (s *twitterServer) GetMentionTimeline(ctx context.Context, req *pb.MentionTimelineRequest) (*pb.Timeline, error) {
  panic("implement me")
}

func (s *twitterServer) GetUserTimeline(ctx context.Context, req *pb.UserTimelineRequest) (*pb.Timeline, error) {
  panic("implement me")
}

func (s *twitterServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.Tweet, error) {
  panic("implement me")
}

func (s *twitterServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.User, error) {
  panic("implement me")
}

func (s *twitterServer) GetRaw(ctx context.Context, rr *pb.RawAPIRequest) (*pb.RawAPIResult, error) {
  auth := decodeAuthPair(rr.Auth)
  or := OAuthRequest{
    Method:   rr.Method,
    Protocol: rr.Protocol,
    Domain:   domain,
    Path:     path.Join(rr.Version, rr.Path),
    Query:    rr.QueryParams,
    Body:     rr.BodyParams,
  }
  req, err := or.MakeRequest(ctx, auth.Secret, auth.Public)
  if err != nil {
    return nil, err
  }
  status, headers, body, err := s.twitter.requestRaw(ctx, req)
  if err != nil {
    return nil, err
  }
  return &pb.RawAPIResult{Status: uint32(status), Headers: headers, Body: body}, nil
}
