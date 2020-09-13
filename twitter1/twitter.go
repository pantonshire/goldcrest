package twitter1

import (
  "encoding/json"
  "fmt"
  "goldcrest"
  "goldcrest/twitter1/model"
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

  xRateLimit          = "X-Rate-Limit-Limit"
  xRateLimitRemaining = "X-Rate-Limit-Remaining"
  xRateLimitReset     = "X-Rate-Limit-Reset"
)

type Config struct {
  ClientTimeoutSeconds uint `json:"client_timeout_seconds"`
}

type Twitter struct {
  client *http.Client
}

type TweetParams struct {
  TrimUser          bool
  IncludeMyRetweet  bool
  IncludeEntities   bool
  IncludeExtAltText bool
  IncludeCardURI    bool
  Mode              TweetMode
}

func NewTwitter(config Config) *Twitter {
  client := &http.Client{
    Timeout: time.Second * time.Duration(config.ClientTimeoutSeconds),
  }

  return &Twitter{
    client: client,
  }
}

func DefaultTweetParams() TweetParams {
  return TweetParams{
    TrimUser:          false,
    IncludeMyRetweet:  true,
    IncludeEntities:   true,
    IncludeExtAltText: true,
    IncludeCardURI:    true,
    Mode:              ExtendedMode,
  }
}

func (t Twitter) requestJSON(req *http.Request, output interface{}) (err error) {
  resp, err := t.client.Do(req)
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
  return json.NewDecoder(resp.Body).Decode(output)
}

func (t Twitter) GetTweet(secret, auth Auth, id interface{}, params TweetParams) (model.Tweet, error) {
  or := OAuthRequest{
    Method:   "GET",
    Protocol: protocol,
    Domain:   domain,
    Path:     path.Join(version, "statuses/show.json"),
    Query: map[string]string{
      "id":                   fmt.Sprint(id),
      "trim_user":            fmt.Sprint(params.TrimUser),
      "include_my_retweet":   fmt.Sprint(params.IncludeMyRetweet),
      "include_entities":     fmt.Sprint(params.IncludeEntities),
      "include_ext_alt_text": fmt.Sprint(params.IncludeExtAltText),
      "include_card_uri":     fmt.Sprint(params.IncludeCardURI),
      "tweet_mode":           string(params.Mode),
    },
  }
  req, err := or.MakeRequest(secret, auth)
  if err != nil {
    return model.Tweet{}, err
  }
  var tweet model.Tweet
  if err := t.requestJSON(req, &tweet); err != nil {
    return model.Tweet{}, err
  }
  return tweet, nil
}
