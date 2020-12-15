package twitter1

import (
  "context"
  "github.com/pantonshire/goldcrest"
  pb "github.com/pantonshire/goldcrest/proto"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "net/http"
  "path"
  "time"
)

type Client interface {
  newContext() (context.Context, context.CancelFunc)
  GetTweet(twOpts TweetOptions, id uint64) (Tweet, error)
  GetHomeTimeline(twOpts TweetOptions, tlOpts TimelineOptions, replies bool) ([]Tweet, error)
  GetMentionTimeline(twOpts TweetOptions, tlOpts TimelineOptions) ([]Tweet, error)
  GetUserIDTimeline(twOpts TweetOptions, id uint64, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error)
  GetUserHandleTimeline(twOpts TweetOptions, handle string, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error)
  UpdateStatus(text string, stOpts StatusUpdateOptions, trimUser bool) (Tweet, error)
  UpdateProfile(pfOpts ProfileUpdateOptions, entities, statuses bool) (User, error)
  GetRaw(method, protocol, version, path string, queryParams, bodyParams map[string]string) (headers map[string]string, status uint, body []byte, err error)
}

func handleRequest(c Client, handler func(ctx context.Context) error) error {
  ctx, cancel := c.newContext()
  if cancel != nil {
    defer cancel()
  }
  err := handler(ctx)
  if httpErr, ok := err.(*goldcrest.HttpError); ok {
    return status.Errorf(codes.Internal, "twitter error %s", httpErr.Error())
  }
  return err
}

type local struct {
  auth    AuthPair
  twitter *Twitter
  timeout time.Duration
}

func Local(auth AuthPair, twitterConfig TwitterConfig, timeout time.Duration) Client {
  return local{
    auth:    auth,
    twitter: NewTwitter(twitterConfig),
    timeout: timeout,
  }
}

func (lc local) newContext() (context.Context, context.CancelFunc) {
  if lc.timeout == 0 {
    return context.Background(), nil
  }
  return context.WithTimeout(context.Background(), lc.timeout)
}

func (lc local) GetTweet(twOpts TweetOptions, id uint64) (Tweet, error) {
  var tweet Tweet
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    tweet, err = lc.twitter.GetTweet(ctx, lc.auth, id, twOpts)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return Tweet{}, err
  }
  return tweet, nil
}

func (lc local) GetHomeTimeline(twOpts TweetOptions, tlOpts TimelineOptions, replies bool) ([]Tweet, error) {
  var tweets []Tweet
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    tweets, err = lc.twitter.GetHomeTimeline(ctx, lc.auth, twOpts, tlOpts, replies)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (lc local) GetMentionTimeline(twOpts TweetOptions, tlOpts TimelineOptions) ([]Tweet, error) {
  var tweets []Tweet
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    tweets, err = lc.twitter.GetMentionTimeline(ctx, lc.auth, twOpts, tlOpts)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (lc local) GetUserIDTimeline(twOpts TweetOptions, id uint64, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  var tweets []Tweet
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    tweets, err = lc.twitter.GetUserTimeline(ctx, lc.auth, twOpts, &id, nil, tlOpts, replies, retweets)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (lc local) GetUserHandleTimeline(twOpts TweetOptions, handle string, tlOpts TimelineOptions, replies, retweets bool) ([]Tweet, error) {
  var tweets []Tweet
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    tweets, err = lc.twitter.GetUserTimeline(ctx, lc.auth, twOpts, nil, &handle, tlOpts, replies, retweets)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return tweets, nil
}

func (lc local) UpdateStatus(text string, stOpts StatusUpdateOptions, trimUser bool) (Tweet, error) {
  var tweet Tweet
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    tweet, err = lc.twitter.UpdateStatus(ctx, lc.auth, text, stOpts.ReplyID, stOpts.AutoReply, stOpts.ExcludeReplyUserIDs, stOpts.AttachmentURL, stOpts.MediaIDs, stOpts.Sensitive, trimUser, stOpts.EnableDMCommands, stOpts.FailDMCommands)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return Tweet{}, err
  }
  return tweet, nil
}

func (lc local) UpdateProfile(pfOpts ProfileUpdateOptions, entities, statuses bool) (User, error) {
  var user User
  err := handleRequest(lc, func(ctx context.Context) error {
    var err error
    user, err = lc.twitter.UpdateProfile(ctx, lc.auth, pfOpts.Name, pfOpts.Url, pfOpts.Location, pfOpts.Bio, pfOpts.LinkColor, entities, statuses)
    if err != nil {
      return err
    }
    return nil
  })
  if err != nil {
    return User{}, err
  }
  return user, nil
}

func (lc local) GetRaw(method, protocol, version, reqPath string, queryParams, bodyParams map[string]string) (headers map[string]string, status uint, body []byte, err error) {
  err = handleRequest(lc, func(ctx context.Context) error {
    or := OAuthRequest{
      Method:   method,
      Protocol: protocol,
      Domain:   domain,
      Path:     path.Join(version, reqPath),
      Query:    queryParams,
      Body:     bodyParams,
    }
    var req *http.Request
    var err error
    req, err = or.MakeRequest(ctx, lc.auth.Secret, lc.auth.Public)
    if err != nil {
      return err
    }
    var statusInt int
    statusInt, headers, body, err = lc.twitter.requestRaw(ctx, req)
    if err != nil {
      return err
    }
    status = uint(statusInt)
    return nil
  })
  if err != nil {
    return nil, 0, nil, err
  }
  return headers, status, body, nil
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

func (rc remote) GetTweet(twOpts TweetOptions, id uint64) (tweet Tweet, err error) {
  err = handleRequest(rc, func(ctx context.Context) error {
    tweetMsg, err := rc.client.GetTweet(ctx, &pb.TweetRequest{
      Auth:    encodeAuthPairMessage(rc.auth),
      Id:      id,
      Options: encodeTweetOptionsMessage(twOpts),
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
  err := handleRequest(rc, func(ctx context.Context) error {
    msg, err := rc.client.GetHomeTimeline(ctx, &pb.HomeTimelineRequest{
      Auth:            encodeAuthPairMessage(rc.auth),
      TimelineOptions: encodeTimelineOptionsMessage(tlOpts),
      TweetOptions:    encodeTweetOptionsMessage(twOpts),
      IncludeReplies:  replies,
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
  err := handleRequest(rc, func(ctx context.Context) error {
    msg, err := rc.client.GetMentionTimeline(ctx, &pb.MentionTimelineRequest{
      Auth:            encodeAuthPairMessage(rc.auth),
      TimelineOptions: encodeTimelineOptionsMessage(tlOpts),
      TweetOptions:    encodeTweetOptionsMessage(twOpts),
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
  err := handleRequest(rc, func(ctx context.Context) error {
    msg, err := rc.client.GetUserTimeline(ctx, &pb.UserTimelineRequest{
      Auth:            encodeAuthPairMessage(rc.auth),
      User:            &pb.UserTimelineRequest_UserId{UserId: id},
      TimelineOptions: encodeTimelineOptionsMessage(tlOpts),
      TweetOptions:    encodeTweetOptionsMessage(twOpts),
      IncludeReplies:  replies,
      IncludeRetweets: retweets,
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
  err := handleRequest(rc, func(ctx context.Context) error {
    msg, err := rc.client.GetUserTimeline(ctx, &pb.UserTimelineRequest{
      Auth:            encodeAuthPairMessage(rc.auth),
      User:            &pb.UserTimelineRequest_UserHandle{UserHandle: handle},
      TimelineOptions: encodeTimelineOptionsMessage(tlOpts),
      TweetOptions:    encodeTweetOptionsMessage(twOpts),
      IncludeReplies:  replies,
      IncludeRetweets: retweets,
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
  err := handleRequest(rc, func(ctx context.Context) error {
    req := pb.UpdateStatusRequest{
      Auth:                      encodeAuthPairMessage(rc.auth),
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
      req.Reply = &pb.UpdateStatusRequest_NoReply{}
    }
    if stOpts.AttachmentURL != nil {
      req.Attachment = &pb.UpdateStatusRequest_AttachmentUrl{AttachmentUrl: *stOpts.AttachmentURL}
    } else {
      req.Attachment = &pb.UpdateStatusRequest_NoAttachment{}
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
  err := handleRequest(rc, func(ctx context.Context) error {
    req := pb.UpdateProfileRequest{
      Auth:            encodeAuthPairMessage(rc.auth),
      IncludeEntities: entities,
      IncludeStatuses: statuses,
    }
    if pfOpts.Name != nil {
      req.UpdateName = &pb.UpdateProfileRequest_Name{Name: *pfOpts.Name}
    } else {
      req.UpdateName = &pb.UpdateProfileRequest_NoUpdateName{}
    }
    if pfOpts.Url != nil {
      req.UpdateUrl = &pb.UpdateProfileRequest_Url{Url: *pfOpts.Url}
    } else {
      req.UpdateUrl = &pb.UpdateProfileRequest_NoUpdateUrl{}
    }
    if pfOpts.Location != nil {
      req.UpdateLocation = &pb.UpdateProfileRequest_Location{Location: *pfOpts.Location}
    } else {
      req.UpdateLocation = &pb.UpdateProfileRequest_NoUpdateLocation{}
    }
    if pfOpts.Bio != nil {
      req.UpdateBio = &pb.UpdateProfileRequest_Bio{Bio: *pfOpts.Bio}
    } else {
      req.UpdateBio = &pb.UpdateProfileRequest_NoUpdateBio{}
    }
    if pfOpts.LinkColor != nil {
      req.UpdateProfileLinkColor = &pb.UpdateProfileRequest_ProfileLinkColor{ProfileLinkColor: *pfOpts.LinkColor}
    } else {
      req.UpdateProfileLinkColor = &pb.UpdateProfileRequest_NoUpdateProfileLinkColor{}
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
  err = handleRequest(rc, func(ctx context.Context) error {
    resp, err := rc.client.GetRaw(ctx, &pb.RawAPIRequest{
      Auth:        encodeAuthPairMessage(rc.auth),
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
