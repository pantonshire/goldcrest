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
  auth := decodeAuthPair(req.Auth)
  twopts := decodeTweetOptions(req.TweetOptions)
  var count *uint
  var minID, maxID *uint64
  if req.Count > 0 {
    countVal := uint(req.Count)
    count = &countVal
  }
  if req.MinId > 0 {
    minID = &req.MinId
  }
  if req.MaxId > 0 {
    maxID = &req.MaxId
  }
  mods, err := s.twitter.GetHomeTimeline(ctx, auth, twopts, count, minID, maxID, req.IncludeReplies)
  if err != nil {
    return nil, err
  }
  return tweetModelsToMessage(mods)
}

func (s *twitterServer) GetMentionTimeline(ctx context.Context, req *pb.MentionTimelineRequest) (*pb.Timeline, error) {
  auth := decodeAuthPair(req.Auth)
  twopts := decodeTweetOptions(req.TweetOptions)
  var count *uint
  var minID, maxID *uint64
  if req.Count > 0 {
    countVal := uint(req.Count)
    count = &countVal
  }
  if req.MinId > 0 {
    minID = &req.MinId
  }
  if req.MaxId > 0 {
    maxID = &req.MaxId
  }
  mods, err := s.twitter.GetMentionTimeline(ctx, auth, twopts, count, minID, maxID)
  if err != nil {
    return nil, err
  }
  return tweetModelsToMessage(mods)
}

func (s *twitterServer) GetUserTimeline(ctx context.Context, req *pb.UserTimelineRequest) (*pb.Timeline, error) {
  auth := decodeAuthPair(req.Auth)
  twopts := decodeTweetOptions(req.TweetOptions)
  var userID *uint64
  var userHandle *string
  var count *uint
  var minID, maxID *uint64
  switch req.User.(type) {
  case *pb.UserTimelineRequest_UserId:
    userID = &req.User.(*pb.UserTimelineRequest_UserId).UserId
  case *pb.UserTimelineRequest_UserHandle:
    userHandle = &req.User.(*pb.UserTimelineRequest_UserHandle).UserHandle
  }
  if req.CountLimit > 0 {
    countVal := uint(req.CountLimit)
    count = &countVal
  }
  if req.MinId > 0 {
    minID = &req.MinId
  }
  if req.MaxId > 0 {
    maxID = &req.MaxId
  }
  mods, err := s.twitter.GetUserTimeline(ctx, auth, twopts, userID, userHandle, count, minID, maxID, req.IncludeReplies, req.IncludeRetweets)
  if err != nil {
    return nil, err
  }
  return tweetModelsToMessage(mods)
}

func (s *twitterServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.Tweet, error) {
  auth := decodeAuthPair(req.Auth)
  var replyID *uint64
  if id, ok := req.Reply.(*pb.UpdateStatusRequest_ReplyId); ok {
    replyID = &id.ReplyId
  }
  var attachmentURL *string
  if url, ok := req.Attachment.(*pb.UpdateStatusRequest_AttachmentUrl); ok {
    attachmentURL = &url.AttachmentUrl
  }
  mod, err := s.twitter.UpdateStatus(ctx, auth, req.Text, replyID, req.AutoPopulateReplyMetadata, req.ExcludeReplyUserIds, attachmentURL, req.MediaIds, req.PossiblySensitive, req.TrimUser, req.EnableDmCommands, req.FailDmCommands)
  if err != nil {
    return nil, err
  }
  return tweetModelToMessage(mod)
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
