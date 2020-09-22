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
  "strconv"
  "sync"
  "time"
)

type TweetMode string
type limitGroup string

const (
  CompatibilityMode TweetMode = "compat"
  ExtendedMode      TweetMode = "extended" //responsible for full_text and extended_entities

  version  = "1.1"
  protocol = "https"
  domain   = "api.twitter.com"

  xRateLimit          = "X-Rate-Limit-Limit"
  xRateLimitRemaining = "X-Rate-Limit-Remaining"
  xRateLimitReset     = "X-Rate-Limit-Reset"

  limitNone         limitGroup = ""
  limitStatusUpdate limitGroup = "statuses/update"
  limitStatusShow   limitGroup = "statuses/show"
)

type TwitterConfig struct {
  ClientTimeoutSeconds uint `json:"client_timeout_seconds"`
}

type Twitter struct {
  client *http.Client
  users  *users
}

type TweetParams struct {
  TrimUser          bool
  IncludeMyRetweet  bool
  IncludeEntities   bool
  IncludeExtAltText bool
  IncludeCardURI    bool
  Mode              TweetMode
}

type users struct {
  lock  sync.Mutex
  cache map[string]*user
}

type user struct {
  lock   sync.Mutex
  limits map[limitGroup]*rateLimit
}

type rateLimit struct {
  lock          sync.Mutex
  current, next int
  resets        time.Time
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

func (us *users) getUser(token string) *user {
  us.lock.Lock()
  defer us.lock.Unlock()
  var u *user
  if u = us.cache[token]; u == nil {
    u = &user{
      limits: make(map[limitGroup]*rateLimit),
    }
    us.cache[token] = u
  }
  return u
}

func (u *user) getRateLimit(group limitGroup) *rateLimit {
  if group == limitNone {
    return nil
  }
  u.lock.Lock()
  defer u.lock.Unlock()
  var rl *rateLimit
  if rl = u.limits[group]; rl == nil {
    rl = &rateLimit{}
    u.limits[group] = rl
  }
  return rl
}

//TODO: retry (with exponential backoff) on 429 too many requests error
func (rl *rateLimit) do(ctx context.Context, fn func() (*http.Response, error)) (*http.Response, error) {
  if rl == nil {
    fmt.Println("Warning: nil rate limit, no block") //for debugging purposes
    return fn()
  }
  rl.lock.Lock()
  defer rl.lock.Unlock()
  if rl.current < 1 {
    if now := time.Now(); rl.resets.After(now) {
      timer := time.NewTimer(rl.resets.Sub(now))
      select {
      case <-timer.C:
      case <-ctx.Done():
        timer.Stop()
        return nil, ctx.Err()
      }
    }
  }
  resp, err := fn()
  if err != nil {
    return resp, err
  }
  if limCurrent := resp.Header.Get(xRateLimitRemaining); limCurrent != "" {
    rl.current, err = strconv.Atoi(limCurrent)
    if err != nil {
      return resp, fmt.Errorf("invalid rate limit header for %s: \"%s\"", xRateLimitRemaining, limCurrent)
    }
  }
  if limNext := resp.Header.Get(xRateLimit); limNext != "" {
    rl.next, err = strconv.Atoi(limNext)
    if err != nil {
      return resp, fmt.Errorf("invalid rate limit header for %s: \"%s\"", xRateLimit, limNext)
    }
  }
  if limResets := resp.Header.Get(xRateLimitReset); limResets != "" {
    resetsUnix, err := strconv.ParseInt(limResets, 10, 64)
    if err != nil {
      return resp, fmt.Errorf("invalid rate limit header for %s: \"%s\"", xRateLimitReset, limResets)
    }
    rl.resets = time.Unix(resetsUnix, 0)
  }
  return resp, nil
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

func (t Twitter) GetTweet(ctx context.Context, secret, auth Auth, id interface{}, params TweetParams) (model.Tweet, error) {
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
  req, err := or.MakeRequest(ctx, secret, auth)
  if err != nil {
    return model.Tweet{}, err
  }
  var tweet model.Tweet
  if err := t.requestJSON(ctx, req, auth.Token, limitStatusShow, &tweet); err != nil {
    return model.Tweet{}, err
  }
  return tweet, nil
}
