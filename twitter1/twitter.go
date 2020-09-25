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

func (t Twitter) GetTweet(ctx context.Context, auth AuthPair, id interface{}, twopts TweetOptions) (model.Tweet, error) {
  or := OAuthRequest{
    Method:   "GET",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/show.json"),
    Query: map[string]string{
      "id":                   fmt.Sprint(id),
      "trim_user":            fmt.Sprint(twopts.TrimUser),
      "include_my_retweet":   fmt.Sprint(twopts.IncludeMyRetweet),
      "include_entities":     fmt.Sprint(twopts.IncludeEntities),
      "include_ext_alt_text": fmt.Sprint(twopts.IncludeExtAltText),
      "include_card_uri":     fmt.Sprint(twopts.IncludeCardURI),
      "tweet_mode":           string(twopts.Mode),
    },
  }
  var tweet model.Tweet
  if err := t.standardRequest(ctx, limitStatusShow, or, auth, &tweet); err != nil {
    return model.Tweet{}, err
  }
  return tweet, nil
}

func (t Twitter) GetHomeTimeline(ctx context.Context, auth AuthPair, twopts TweetOptions, count *uint, minID, maxID *uint64, includeReplies bool) ([]model.Tweet, error) {
  query := map[string]string{
    "trim_user":            fmt.Sprint(twopts.TrimUser),
    "include_my_retweet":   fmt.Sprint(twopts.IncludeMyRetweet),
    "include_entities":     fmt.Sprint(twopts.IncludeEntities),
    "include_ext_alt_text": fmt.Sprint(twopts.IncludeExtAltText),
    "include_card_uri":     fmt.Sprint(twopts.IncludeCardURI),
    "tweet_mode":           string(twopts.Mode),
    "exclude_replies":      fmt.Sprint(!includeReplies),
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
    Path:     path.Join(version, "statuses/home_timeline.json"),
    Query:    query,
  }
  var tweets []model.Tweet
  if err := t.standardRequest(ctx, limitHomeTimeline, or, auth, &tweets); err != nil {
    return nil, err
  }
  return tweets, nil
}
