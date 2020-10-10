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
  GetTweet(twOpts TweetOptions, id uint64) (Tweet, error)
  GetHomeTimeline(twOpts TweetOptions, tlOpts TimelineOptions, replies bool) ([]Tweet, error)
  GetMentionTimeline(twOpts TweetOptions, tlOpts TimelineOptions) ([]Tweet, error)
  GetUserIDTimeline(twOpts TweetOptions, id uint64, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error)
  GetUserHandleTimeline(twOpts TweetOptions, handle string, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error)
  UpdateStatus(text string, stOpts StatusUpdateOptions, trimUser bool) (Tweet, error)
  UpdateProfile(pfOpts ProfileUpdateOptions, entities, statuses bool) (User, error)
  GetRaw(method, protocol, version, path string, queryParams, bodyParams map[string]string) (headers map[string]string, status uint, body []byte, err error)
}

type local struct {
  auth    AuthPair
  twitter *Twitter
}

func Local(auth AuthPair, twitterConfig TwitterConfig) Client {
  return local{
    auth:    auth,
    twitter: NewTwitter(twitterConfig),
  }
}

func (lc local) GetTweet(twOpts TweetOptions, id uint64) (Tweet, error) {
  panic("implement me")
}

func (lc local) GetHomeTimeline(twOpts TweetOptions, tlOpts TimelineOptions, replies bool) ([]Tweet, error) {
  panic("implement me")
}

func (lc local) GetMentionTimeline(twOpts TweetOptions, tlOpts TimelineOptions) ([]Tweet, error) {
  panic("implement me")
}

func (lc local) GetUserIDTimeline(twOpts TweetOptions, id uint64, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  panic("implement me")
}

func (lc local) GetUserHandleTimeline(twOpts TweetOptions, handle string, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  panic("implement me")
}

func (lc local) UpdateStatus(text string, stOpts StatusUpdateOptions, trimUser bool) (Tweet, error) {
  panic("implement me")
}

func (lc local) UpdateProfile(pfOpts ProfileUpdateOptions, entities, statuses bool) (User, error) {
  panic("implement me")
}

func (lc local) GetRaw(method, protocol, version, path string, queryParams, bodyParams map[string]string) (headers map[string]string, status uint, body []byte, err error) {
  panic("implement me")
}

type remote struct {
  auth        AuthPair
  address     string
  client      pb.Twitter1Client
  callTimeout time.Duration
}

func Remote(conn *grpc.ClientConn, auth AuthPair, timeout time.Duration) Client {
  return remote{
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

func (rc remote) handleRequest(handler func(ctx context.Context) error) error {
  ctx, cancel := rc.newContext()
  if cancel != nil {
    defer cancel()
  }
  err := handler(ctx)
  if httpErr, ok := err.(*goldcrest.HttpError); ok {
    return status.Errorf(codes.Internal, "twitter error %s", httpErr.Error())
  }
  return err
}

func (rc remote) GetTweet(twOpts TweetOptions, id uint64) (tweet Tweet, err error) {
  err = rc.handleRequest(func(ctx context.Context) error {
    tweetMsg, err := rc.client.GetTweet(ctx, &pb.TweetRequest{
      Auth:    authPairToMsg(rc.auth),
      Id:      id,
      Options: tweetOptionsToMsg(twOpts),
    })
    if err != nil {
      return err
    }
    tweet = decodeTweetMessage(tweetMsg)
    return nil
  })
  if err != nil {
    return Tweet{}, err
  }
  return tweet, nil
}

func (rc remote) GetHomeTimeline(twOpts TweetOptions, tlOpts TimelineOptions, replies bool) ([]Tweet, error) {
  var tweets []Tweet
  err := rc.handleRequest(func(ctx context.Context) error {
    msg, err := rc.client.GetHomeTimeline(ctx, &pb.HomeTimelineRequest{
      Auth:           authPairToMsg(rc.auth),
      Count:          uint32(tlOpts.Count),
      MinId:          tlOpts.MinID,
      MaxId:          tlOpts.MaxID,
      IncludeReplies: replies,
      TweetOptions:   tweetOptionsToMsg(twOpts),
    })
    if err != nil {
      return err
    }
    tweets = decodeTimelineMessage(msg)
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (rc remote) GetMentionTimeline(twOpts TweetOptions, tlOpts TimelineOptions) ([]Tweet, error) {
  var tweets []Tweet
  err := rc.handleRequest(func(ctx context.Context) error {
    msg, err := rc.client.GetMentionTimeline(ctx, &pb.MentionTimelineRequest{
      Auth:         authPairToMsg(rc.auth),
      Count:        uint32(tlOpts.Count),
      MinId:        tlOpts.MinID,
      MaxId:        tlOpts.MaxID,
      TweetOptions: tweetOptionsToMsg(twOpts),
    })
    if err != nil {
      return err
    }
    tweets = decodeTimelineMessage(msg)
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (rc remote) GetUserIDTimeline(twOpts TweetOptions, id uint64, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  var tweets []Tweet
  err := rc.handleRequest(func(ctx context.Context) error {
    msg, err := rc.client.GetUserTimeline(ctx, &pb.UserTimelineRequest{
      Auth:            authPairToMsg(rc.auth),
      User:            &pb.UserTimelineRequest_UserId{UserId: id},
      CountLimit:      uint32(tlOpts.Count),
      MinId:           tlOpts.MinID,
      MaxId:           tlOpts.MaxID,
      IncludeReplies:  replies,
      IncludeRetweets: retweets,
      TweetOptions:    tweetOptionsToMsg(twOpts),
    })
    if err != nil {
      return err
    }
    tweets = decodeTimelineMessage(msg)
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (rc remote) GetUserHandleTimeline(twOpts TweetOptions, handle string, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  var tweets []Tweet
  err := rc.handleRequest(func(ctx context.Context) error {
    msg, err := rc.client.GetUserTimeline(ctx, &pb.UserTimelineRequest{
      Auth:            authPairToMsg(rc.auth),
      User:            &pb.UserTimelineRequest_UserHandle{UserHandle: handle},
      CountLimit:      uint32(tlOpts.Count),
      MinId:           tlOpts.MinID,
      MaxId:           tlOpts.MaxID,
      IncludeReplies:  replies,
      IncludeRetweets: retweets,
      TweetOptions:    tweetOptionsToMsg(twOpts),
    })
    if err != nil {
      return err
    }
    tweets = decodeTimelineMessage(msg)
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (rc remote) UpdateStatus(text string, stOpts StatusUpdateOptions, trimUser bool) (Tweet, error) {
  var tweet Tweet
  err := rc.handleRequest(func(ctx context.Context) error {
    req := pb.UpdateStatusRequest{
      Auth:                      authPairToMsg(rc.auth),
      Text:                      text,
      AutoPopulateReplyMetadata: stOpts.AutoReply,
      ExcludeReplyUserIds:       stOpts.ExcludeReplyUserIDs,
      MediaIds:                  stOpts.MediaIDs,
      PossiblySensitive:         stOpts.Sensitive,
      TrimUser:                  trimUser,
      EnableDmCommands:          stOpts.EnableDMCommands,
      FailDmCommands:            stOpts.FailDMCommands,
    }
    if stOpts.ReplyID != nil {
      req.Reply = &pb.UpdateStatusRequest_ReplyId{ReplyId: *stOpts.ReplyID}
    } else {
      req.Reply = &pb.UpdateStatusRequest_NoReply{NoReply: true}
    }
    if stOpts.AttachmentURL != nil {
      req.Attachment = &pb.UpdateStatusRequest_AttachmentUrl{AttachmentUrl: *stOpts.AttachmentURL}
    } else {
      req.Attachment = &pb.UpdateStatusRequest_NoAttachment{NoAttachment: true}
    }
    msg, err := rc.client.UpdateStatus(ctx, &req)
    if err != nil {
      return err
    }
    tweet = decodeTweetMessage(msg)
    return nil
  })
  if err != nil {
    return Tweet{}, err
  }
  return tweet, nil
}

func (rc remote) UpdateProfile(pfOpts ProfileUpdateOptions, entities, statuses bool) (User, error) {
  var user User
  err := rc.handleRequest(func(ctx context.Context) error {
    req := pb.UpdateProfileRequest{
      Auth:            authPairToMsg(rc.auth),
      IncludeEntities: entities,
      IncludeStatuses: statuses,
    }
    if pfOpts.Name != nil {
      req.UpdateName = &pb.UpdateProfileRequest_Name{Name: *pfOpts.Name}
    } else {
      req.UpdateName = &pb.UpdateProfileRequest_NoUpdateName{NoUpdateName: true}
    }
    if pfOpts.Url != nil {
      req.UpdateUrl = &pb.UpdateProfileRequest_Url{Url: *pfOpts.Url}
    } else {
      req.UpdateUrl = &pb.UpdateProfileRequest_NoUpdateUrl{NoUpdateUrl: true}
    }
    if pfOpts.Location != nil {
      req.UpdateLocation = &pb.UpdateProfileRequest_Location{Location: *pfOpts.Location}
    } else {
      req.UpdateLocation = &pb.UpdateProfileRequest_NoUpdateLocation{NoUpdateLocation: true}
    }
    if pfOpts.Bio != nil {
      req.UpdateBio = &pb.UpdateProfileRequest_Bio{Bio: *pfOpts.Bio}
    } else {
      req.UpdateBio = &pb.UpdateProfileRequest_NoUpdateBio{NoUpdateBio: true}
    }
    if pfOpts.LinkColor != nil {
      req.UpdateProfileLinkColor = &pb.UpdateProfileRequest_ProfileLinkColor{ProfileLinkColor: *pfOpts.LinkColor}
    } else {
      req.UpdateProfileLinkColor = &pb.UpdateProfileRequest_NoUpdateProfileLinkColor{NoUpdateProfileLinkColor: true}
    }
    msg, err := rc.client.UpdateProfile(ctx, &req)
    if err != nil {
      return err
    }
    user = decodeUserMessage(msg)
    return nil
  })
  if err != nil {
    return User{}, err
  }
  return user, nil
}

func (rc remote) GetRaw(method, protocol, version, path string, queryParams, bodyParams map[string]string) (headers map[string]string, status uint, body []byte, err error) {
  err = rc.handleRequest(func(ctx context.Context) error {
    resp, err := rc.client.GetRaw(ctx, &pb.RawAPIRequest{
      Auth:        authPairToMsg(rc.auth),
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
