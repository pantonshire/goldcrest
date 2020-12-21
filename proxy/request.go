package proxy

import (
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/oauth"
  "strconv"
  "strings"
)

type tweetOptions struct {
  trimUser          bool
  includeMyRetweet  bool
  includeEntities   bool
  includeExtAltText bool
  includeCardURI    bool
  mode              tweetMode
}

type timelineOptions struct {
  count     uint
  minID     *uint64
  maxID     *uint64
  tweetOpts tweetOptions
}

var (
  defaultTweetOptions = tweetOptions{
    trimUser:          false,
    includeMyRetweet:  true,
    includeEntities:   true,
    includeExtAltText: true,
    includeCardURI:    true,
    mode:              extendedMode,
  }
  defaultTweetOptionsParams = defaultTweetOptions.ser()
)

func desAuth(msg *pb.Authentication) oauth.AuthPair {
  if msg == nil {
    return oauth.AuthPair{}
  }
  return oauth.AuthPair{
    Secret: oauth.Auth{
      Key:   msg.SecretKey,
      Token: msg.SecretToken,
    },
    Public: oauth.Auth{
      Key:   msg.ConsumerKey,
      Token: msg.AccessToken,
    },
  }
}

func desTweetMode(msg pb.TweetOptions_Mode) tweetMode {
  if msg == pb.TweetOptions_EXTENDED {
    return extendedMode
  }
  return compatibilityMode
}

func desTweetOptions(options *pb.TweetOptions) tweetOptions {
  if options == nil {
    return tweetOptions{}
  }
  return tweetOptions{
    trimUser:          options.TrimUser,
    includeMyRetweet:  options.IncludeMyRetweet,
    includeEntities:   options.IncludeEntities,
    includeExtAltText: options.IncludeExtAltText,
    includeCardURI:    options.IncludeCardUri,
    mode:              desTweetMode(options.Mode),
  }
}

func (opts tweetOptions) ser() oauth.Params {
  return map[string]string{
    "trim_user":            strconv.FormatBool(opts.trimUser),
    "include_my_retweet":   strconv.FormatBool(opts.includeMyRetweet),
    "include_entities":     strconv.FormatBool(opts.includeEntities),
    "include_ext_alt_text": strconv.FormatBool(opts.includeExtAltText),
    "include_card_uri":     strconv.FormatBool(opts.includeCardURI),
    "tweet_mode":           opts.mode.String(),
  }
}

func desTimelineOptions(msg *pb.TimelineOptions) timelineOptions {
  if msg == nil {
    return timelineOptions{}
  }
  opts := timelineOptions{
    count: uint(msg.Count),
  }
  if minID, ok := msg.Min.(*pb.TimelineOptions_MinId); ok {
    opts.minID = new(uint64)
    *opts.minID = minID.MinId
  }
  if maxID, ok := msg.Max.(*pb.TimelineOptions_MaxId); ok {
    opts.maxID = new(uint64)
    *opts.maxID = maxID.MaxId
  }
  if tweetOpts, ok := msg.Content.(*pb.TimelineOptions_Custom); ok {
    opts.tweetOpts = desTweetOptions(tweetOpts.Custom)
  } else {
    opts.tweetOpts = defaultTweetOptions
  }
  return opts
}

func (opts timelineOptions) ser() oauth.Params {
  params := oauth.NewParams()
  params.Extend(opts.tweetOpts.ser())
  params.Set("count", strconv.FormatUint(uint64(opts.count), 10))
  if opts.minID != nil && *opts.minID > 0 {
    params.Set("since_id", strconv.FormatUint(*opts.minID-1, 10))
  }
  if opts.maxID != nil {
    params.Set("max_id", strconv.FormatUint(*opts.maxID, 10))
  }
  return params
}

func reserTweetRequest(msg *pb.TweetRequest) (oauth.AuthPair, oauth.Params) {
  if msg == nil {
    return oauth.AuthPair{}, nil
  }
  auth := desAuth(msg.Auth)
  params := oauth.NewParams()
  params.Set("id", strconv.FormatUint(msg.Id, 10))
  if custom, ok := msg.Content.(*pb.TweetRequest_Custom); ok {
    params.Extend(desTweetOptions(custom.Custom).ser())
  } else {
    params.Extend(defaultTweetOptionsParams)
  }
  return auth, params
}

func reserTweetsRequest(msg *pb.TweetsRequest) (oauth.AuthPair, oauth.Params) {
  if msg == nil {
    return oauth.AuthPair{}, nil
  }
  auth := desAuth(msg.Auth)
  params := oauth.NewParams()
  if len(msg.Ids) > 0 {
    ids := make([]string, len(msg.Ids))
    for i, id := range msg.Ids {
      ids[i] = strconv.FormatUint(id, 10)
    }
    params.Set("id", strings.Join(ids, ","))
  }
  if custom, ok := msg.Content.(*pb.TweetsRequest_Custom); ok {
    params.Extend(desTweetOptions(custom.Custom).ser())
  } else {
    params.Extend(defaultTweetOptionsParams)
  }
  return auth, params
}

func reserPublishTweetRequest(msg *pb.PublishTweetRequest) (oauth.AuthPair, oauth.Params) {
  if msg == nil {
    return oauth.AuthPair{}, nil
  }
  auth := desAuth(msg.Auth)
  params := oauth.NewParams()
  params.Set("status", msg.Text)
  params.Set("auto_populate_reply_metadata", strconv.FormatBool(msg.AutoPopulateReplyMetadata))
  params.Set("possibly_sensitive", strconv.FormatBool(msg.PossiblySensitive))
  params.Set("enable_dmcommands", strconv.FormatBool(msg.EnableDmCommands))
  params.Set("fail_dmcommands", strconv.FormatBool(msg.FailDmCommands))
  if custom, ok := msg.Content.(*pb.PublishTweetRequest_Custom); ok {
    params.Extend(desTweetOptions(custom.Custom).ser())
  } else {
    params.Extend(defaultTweetOptionsParams)
  }
  return auth, params
}
