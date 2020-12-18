package proxy

import (
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/oauth"
  "strconv"
)

type tweetOptions struct {
  trimUser          bool
  includeMyRetweet  bool
  includeEntities   bool
  includeExtAltText bool
  includeCardURI    bool
  mode              tweetMode
}

func (opts tweetOptions) encode() map[string]string {
  return map[string]string{
    "trim_user":            strconv.FormatBool(opts.trimUser),
    "include_my_retweet":   strconv.FormatBool(opts.includeMyRetweet),
    "include_entities":     strconv.FormatBool(opts.includeEntities),
    "include_ext_alt_text": strconv.FormatBool(opts.includeExtAltText),
    "include_card_uri":     strconv.FormatBool(opts.includeCardURI),
    "tweet_mode":           opts.mode.String(),
  }
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
)

func desTweetRequest(msg *pb.TweetRequest) (oauth.AuthPair, uint64, tweetOptions) {
  if msg == nil {
    return oauth.AuthPair{}, 0, tweetOptions{}
  }
  auth := desAuth(msg.Auth)
  id := msg.Id
  var tweetOpts tweetOptions
  if custom, ok := msg.Content.(*pb.TweetRequest_Custom); ok {
    tweetOpts = desTweetOptions(custom.Custom)
  } else {
    tweetOpts = defaultTweetOptions
  }
  return auth, id, tweetOpts
}

func desTweetsRequest(msg *pb.TweetsRequest) (oauth.AuthPair, []uint64, tweetOptions) {
  if msg == nil {
    return oauth.AuthPair{}, nil, tweetOptions{}
  }
  auth := desAuth(msg.Auth)
  ids := make([]uint64, len(msg.Ids))
  copy(ids, msg.Ids)
  var tweetOpts tweetOptions
  if custom, ok := msg.Content.(*pb.TweetsRequest_Custom); ok {
    tweetOpts = desTweetOptions(custom.Custom)
  } else {
    tweetOpts = defaultTweetOptions
  }
  return auth, ids, tweetOpts
}

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
