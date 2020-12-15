package twitter1

import (
  "context"
  pb "github.com/pantonshire/goldcrest/proto"
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

func (s *twitterServer) GetTweet(ctx context.Context, req *pb.TweetRequest) (*pb.TweetResponse, error) {
  auth := decodeAuthPairMessage(req.Auth)
  var opts TweetOptions
  if customOpts, ok := req.Content.(*pb.TweetRequest_Custom); ok {
    opts = decodeTweetOptionsMessage(customOpts.Custom)
  } else {
    opts = DefaultTweetOptions()
  }
  tweet, err := s.twitter.GetTweet(ctx, auth, req.Id, opts)
  if err != nil {
    return nil, err
  }
  return encodeTweetMessage(tweet), nil
}

func (s *twitterServer) GetHomeTimeline(ctx context.Context, req *pb.HomeTimelineRequest) (*pb.TimelineResponse, error) {
  auth := decodeAuthPairMessage(req.Auth)
  twOpts := decodeTweetOptionsMessage(req.TweetOptions)
  tlOpts := decodeTimelineOptionsMessage(req.TimelineOptions)
  tweets, err := s.twitter.GetHomeTimeline(ctx, auth, twOpts, tlOpts, req.IncludeReplies)
  if err != nil {
    return nil, err
  }
  return encodeTimelineMessage(tweets), nil
}

func (s *twitterServer) GetMentionTimeline(ctx context.Context, req *pb.MentionTimelineRequest) (*pb.TimelineResponse, error) {
  auth := decodeAuthPairMessage(req.Auth)
  twOpts := decodeTweetOptionsMessage(req.TweetOptions)
  tlOpts := decodeTimelineOptionsMessage(req.TimelineOptions)
  tweets, err := s.twitter.GetMentionTimeline(ctx, auth, twOpts, tlOpts)
  if err != nil {
    return nil, err
  }
  return encodeTimelineMessage(tweets), nil
}

func (s *twitterServer) GetUserTimeline(ctx context.Context, req *pb.UserTimelineRequest) (*pb.TimelineResponse, error) {
  auth := decodeAuthPairMessage(req.Auth)
  twOpts := decodeTweetOptionsMessage(req.TweetOptions)
  tlOpts := decodeTimelineOptionsMessage(req.TimelineOptions)
  var userID *uint64
  var userHandle *string
  switch req.User.(type) {
  case *pb.UserTimelineRequest_UserId:
    userID = &req.User.(*pb.UserTimelineRequest_UserId).UserId
  case *pb.UserTimelineRequest_UserHandle:
    userHandle = &req.User.(*pb.UserTimelineRequest_UserHandle).UserHandle
  }
  tweets, err := s.twitter.GetUserTimeline(ctx, auth, twOpts, userID, userHandle, tlOpts, req.IncludeReplies, req.IncludeRetweets)
  if err != nil {
    return nil, err
  }
  return encodeTimelineMessage(tweets), nil
}

func (s *twitterServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.TweetResponse, error) {
  auth := decodeAuthPairMessage(req.Auth)
  var replyID *uint64
  if id, ok := req.Reply.(*pb.UpdateStatusRequest_ReplyId); ok {
    replyID = &id.ReplyId
  }
  var attachmentURL *string
  if url, ok := req.Attachment.(*pb.UpdateStatusRequest_AttachmentUrl); ok {
    attachmentURL = &url.AttachmentUrl
  }
  tweet, err := s.twitter.UpdateStatus(ctx, auth, req.Text, replyID, req.AutoPopulateReplyMetadata, req.ExcludeReplyUserIds, attachmentURL, req.MediaIds, req.PossiblySensitive, req.TrimUser, req.EnableDmCommands, req.FailDmCommands)
  if err != nil {
    return nil, err
  }
  return encodeTweetMessage(tweet), nil
}

func (s *twitterServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
  auth := decodeAuthPairMessage(req.Auth)
  var name *string
  if val, ok := req.UpdateName.(*pb.UpdateProfileRequest_Name); ok {
    name = &val.Name
  }
  var url *string
  if val, ok := req.UpdateUrl.(*pb.UpdateProfileRequest_Url); ok {
    url = &val.Url
  }
  var location *string
  if val, ok := req.UpdateLocation.(*pb.UpdateProfileRequest_Location); ok {
    location = &val.Location
  }
  var bio *string
  if val, ok := req.UpdateBio.(*pb.UpdateProfileRequest_Bio); ok {
    bio = &val.Bio
  }
  var linkColor *string
  if val, ok := req.UpdateProfileLinkColor.(*pb.UpdateProfileRequest_ProfileLinkColor); ok {
    linkColor = &val.ProfileLinkColor
  }
  user, err := s.twitter.UpdateProfile(ctx, auth, name, url, location, bio, linkColor, req.IncludeEntities, req.IncludeStatuses)
  if err != nil {
    return nil, err
  }
  return encodeUserMessage(user), nil
}

func (s *twitterServer) GetRaw(ctx context.Context, rr *pb.RawAPIRequest) (*pb.RawAPIResult, error) {
  auth := decodeAuthPairMessage(rr.Auth)
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
