package twitter1

import (
  "context"
  "encoding/json"
  "fmt"
  "goldcrest"
  "goldcrest/twitter1/model"
  "io/ioutil"
  "net/http"
  "path"
  "strings"
  "time"
)

type TweetMode string

const (
  CompatibilityMode TweetMode = "compat"
  ExtendedMode      TweetMode = "extended" //responsible for full_text and extended_entities

  version  = "1.1"
  protocol = "https"
  domain   = "api.twitter.com"
)

type TwitterConfig struct {
  ClientTimeoutSeconds uint `json:"client_timeout_seconds"`
}

type Twitter struct {
  client *http.Client
  users  *users
}

type TweetOptions struct {
  TrimUser          bool
  IncludeMyRetweet  bool
  IncludeEntities   bool
  IncludeExtAltText bool
  IncludeCardURI    bool
  Mode              TweetMode
}

type TimelineOptions struct {
  count        uint
  minID, maxID uint64
}

type StatusUpdateOptions struct {
  replyID                          *uint64
  autoReply                        bool
  excludeReplyUserIDs              []uint64
  attachmentURL                    *string
  mediaIDs                         []uint64
  sensitive                        bool
  trimUser                         bool
  enableDMCommands, failDMCommands bool
}

type ProfileUpdateOptions struct {
  name      *string
  url       *string
  location  *string
  bio       *string
  linkColor *string
}

func NewTwitter(config TwitterConfig) *Twitter {
  client := &http.Client{
    Timeout: time.Second * time.Duration(config.ClientTimeoutSeconds),
  }
  return &Twitter{
    client: client,
    users:  &users{cache: make(map[string]*user)},
  }
}

func DefaultTweetOptions() TweetOptions {
  return TweetOptions{
    TrimUser:          false,
    IncludeMyRetweet:  true,
    IncludeEntities:   true,
    IncludeExtAltText: true,
    IncludeCardURI:    true,
    Mode:              ExtendedMode,
  }
}

func (twOpts TweetOptions) encode() map[string]string {
  return map[string]string{
    "trim_user":            fmt.Sprint(twOpts.TrimUser),
    "include_my_retweet":   fmt.Sprint(twOpts.IncludeMyRetweet),
    "include_entities":     fmt.Sprint(twOpts.IncludeEntities),
    "include_ext_alt_text": fmt.Sprint(twOpts.IncludeExtAltText),
    "include_card_uri":     fmt.Sprint(twOpts.IncludeCardURI),
    "tweet_mode":           string(twOpts.Mode),
  }
}

func (t Twitter) request(ctx context.Context, req *http.Request, token string, group limitGroup, handler func(resp *http.Response) error) (err error) {
  var resp *http.Response
  if token != "" {
    resp, err = t.users.getUser(token).getRateLimit(group).do(ctx, func() (*http.Response, error) {
      return t.client.Do(req)
    })
  } else {
    resp, err = t.client.Do(req)
  }
  if err != nil {
    return err
  }
  defer func() {
    if closeErr := resp.Body.Close(); closeErr != nil {
      err = closeErr
    }
  }()
  httpErr := goldcrest.HttpErrorFor(resp)
  if httpErr != nil {
    return httpErr
  }
  return handler(resp)
}

func (t Twitter) requestJSON(ctx context.Context, req *http.Request, token string, group limitGroup, output interface{}) (err error) {
  return t.request(ctx, req, token, group, func(resp *http.Response) error {
    return json.NewDecoder(resp.Body).Decode(output)
  })
}

func (t Twitter) standardRequest(ctx context.Context, group limitGroup, or OAuthRequest, auth AuthPair, output interface{}) error {
  req, err := or.MakeRequest(ctx, auth.Secret, auth.Public)
  if err != nil {
    return err
  }
  if err := t.requestJSON(ctx, req, auth.Public.Token, group, output); err != nil {
    return err
  }
  return nil
}

func (t Twitter) requestRaw(ctx context.Context, req *http.Request) (status int, headers map[string]string, body []byte, err error) {
  err = t.request(ctx, req, "", limitNone, func(resp *http.Response) error {
    headers = make(map[string]string)
    for key, val := range resp.Header {
      if len(val) > 0 {
        headers[key] = val[0]
      }
    }
    var ioErr error
    body, ioErr = ioutil.ReadAll(resp.Body)
    if err != nil {
      return ioErr
    }
    status = resp.StatusCode
    return nil
  })
  if err != nil {
    return 0, nil, nil, err
  }
  return status, headers, body, nil
}

func (t Twitter) GetTweet(ctx context.Context, auth AuthPair, id interface{}, twOpts TweetOptions) (model.Tweet, error) {
  or := OAuthRequest{
    Method:   "GET",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/show.json"),
    Query: joinParamMaps(map[string]string{
      "id": fmt.Sprint(id),
    }, twOpts.encode()),
  }
  var tweet model.Tweet
  if err := t.standardRequest(ctx, limitStatusShow, or, auth, &tweet); err != nil {
    return model.Tweet{}, err
  }
  return tweet, nil
}

func (t Twitter) GetHomeTimeline(ctx context.Context, auth AuthPair, twOpts TweetOptions, count *uint, minID, maxID *uint64, includeReplies bool) ([]model.Tweet, error) {
  query := joinParamMaps(map[string]string{
    "exclude_replies": fmt.Sprint(!includeReplies),
  }, twOpts.encode())
  if count != nil {
    query["count"] = fmt.Sprint(*count)
  }
  if minID != nil {
    query["since_id"] = fmt.Sprint(*minID)
  }
  if maxID != nil {
    query["max_id"] = fmt.Sprint(*maxID)
  }
  or := OAuthRequest{
    Method:   "GET",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/home_timeline.json"),
    Query:    query,
  }
  var tweets []model.Tweet
  if err := t.standardRequest(ctx, limitHomeTimeline, or, auth, &tweets); err != nil {
    return nil, err
  }
  return tweets, nil
}

func (t Twitter) GetMentionTimeline(ctx context.Context, auth AuthPair, twOpts TweetOptions, count *uint, minID, maxID *uint64) ([]model.Tweet, error) {
  query := twOpts.encode()
  if count != nil {
    query["count"] = fmt.Sprint(*count)
  }
  if minID != nil {
    query["since_id"] = fmt.Sprint(*minID)
  }
  if maxID != nil {
    query["max_id"] = fmt.Sprint(*maxID)
  }
  or := OAuthRequest{
    Method:   "GET",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/mentions_timeline.json"),
    Query:    query,
  }
  var tweets []model.Tweet
  if err := t.standardRequest(ctx, limitMentionTimeline, or, auth, &tweets); err != nil {
    return nil, err
  }
  return tweets, nil
}

func (t Twitter) GetUserTimeline(ctx context.Context, auth AuthPair, twOpts TweetOptions, id *uint64, handle *string, count *uint, minID, maxID *uint64, includeReplies, includeRetweets bool) ([]model.Tweet, error) {
  query := joinParamMaps(map[string]string{
    "exclude_replies": fmt.Sprint(!includeReplies),
    "include_rts":     fmt.Sprint(includeRetweets),
  }, twOpts.encode())
  if id != nil {
    query["user_id"] = fmt.Sprint(*id)
  }
  if handle != nil {
    query["screen_name"] = fmt.Sprint(*handle)
  }
  if count != nil {
    query["count"] = fmt.Sprint(*count)
  }
  if minID != nil {
    query["since_id"] = fmt.Sprint(*minID)
  }
  if maxID != nil {
    query["max_id"] = fmt.Sprint(*maxID)
  }
  or := OAuthRequest{
    Method:   "GET",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/user_timeline.json"),
    Query:    query,
  }
  var tweets []model.Tweet
  if err := t.standardRequest(ctx, limitUserTimeline, or, auth, &tweets); err != nil {
    return nil, err
  }
  return tweets, nil
}

func (t Twitter) UpdateStatus(ctx context.Context, auth AuthPair, text string, replyID *uint64, autoReply bool, excludeReplyUserIDs []uint64, attachmentURL *string, mediaIDs []uint64, sensitive, trimUser, enableDMCommands, failDMCommands bool) (model.Tweet, error) {
  query := map[string]string{
    "status":                       text,
    "auto_populate_reply_metadata": fmt.Sprint(autoReply),
    "possibly_sensitive":           fmt.Sprint(sensitive),
    "trim_user":                    fmt.Sprint(trimUser),
    "enable_dmcommands":            fmt.Sprint(enableDMCommands),
    "fail_dmcommands":              fmt.Sprint(failDMCommands),
  }
  if replyID != nil {
    query["in_reply_to_status_id"] = fmt.Sprint(*replyID)
  }
  if attachmentURL != nil {
    query["attachment_url"] = fmt.Sprint(*attachmentURL)
  }
  if len(excludeReplyUserIDs) > 0 {
    strs := make([]string, len(excludeReplyUserIDs))
    for i, id := range excludeReplyUserIDs {
      strs[i] = fmt.Sprint(id)
    }
    query["exclude_reply_user_ids"] = strings.Join(strs, ",")
  }
  if len(mediaIDs) > 0 {
    strs := make([]string, len(mediaIDs))
    for i, id := range mediaIDs {
      strs[i] = fmt.Sprint(id)
    }
    query["media_ids"] = strings.Join(strs, ",")
  }
  or := OAuthRequest{
    Method:   "POST",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/update.json"),
    Query:    query,
  }
  var tweet model.Tweet
  if err := t.standardRequest(ctx, limitStatusUpdate, or, auth, &tweet); err != nil {
    return model.Tweet{}, err
  }
  return tweet, nil
}

func (t Twitter) UpdateProfile(ctx context.Context, auth AuthPair, name, url, location, bio, linkColor *string, includeEntities, includeStatuses bool) (model.User, error) {
  query := map[string]string{
    "include_entities": fmt.Sprint(includeEntities),
    "skip_status":      fmt.Sprint(!includeStatuses),
  }
  if name != nil {
    query["name"] = *name
  }
  if url != nil {
    query["url"] = *url
  }
  if location != nil {
    query["location"] = *location
  }
  if bio != nil {
    query["description"] = *bio
  }
  if linkColor != nil {
    query["profile_link_color"] = *linkColor
  }
  or := OAuthRequest{
    Method:   "POST",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/update.json"),
    Query:    query,
  }
  var user model.User
  if err := t.standardRequest(ctx, limitUpdateProfile, or, auth, &user); err != nil {
    return model.User{}, err
  }
  return user, nil
}

func joinParamMaps(ms ...map[string]string) map[string]string {
  master := make(map[string]string)
  for _, m := range ms {
    for key, val := range m {
      master[key] = val
    }
  }
  return master
}
