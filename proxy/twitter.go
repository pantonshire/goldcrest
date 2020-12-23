package proxy

import (
  "encoding/json"
  "fmt"
  "github.com/pantonshire/goldcrest/proxy/oauth"
  "io/ioutil"
  "math/bits"
  "net/http"
  "strconv"
  "time"
)

const (
  headerRateLimit          = "X-Rate-Limit-Limit"
  headerRateLimitRemaining = "X-Rate-Limit-Remaining"
  headerRateLimitReset     = "X-Rate-Limit-Reset"
)

type tweetMode string

const (
  compatibilityMode tweetMode = "compat"
  extendedMode      tweetMode = "extended"
)

func (tm tweetMode) String() string {
  return string(tm)
}

type limitGroup string

const (
  publishLimitGroup = "publish"
)

type requestMethod string

const (
  methodGet  requestMethod = "GET"
  methodPost requestMethod = "POST"
)

func (method requestMethod) String() string {
  return string(method)
}

type endpoint struct {
  path   string
  method requestMethod
  group  limitGroup
}

var (
  showTweetEndpoint       = endpoint{path: "statuses/show.json", method: methodGet}
  showTweetsEndpoint      = endpoint{path: "statuses/lookup.json", method: methodGet}
  homeTimelineEndpoint    = endpoint{path: "statuses/home_timeline.json", method: methodGet}
  mentionTimelineEndpoint = endpoint{path: "statuses/mentions_timeline.json", method: methodGet}
  userTimelineEndpoint    = endpoint{path: "statuses/user_timeline.json", method: methodGet}
  publishTweetEndpoint    = endpoint{path: "statuses/update.json", method: methodPost, group: publishLimitGroup}
  destroyTweetEndpoint    = endpoint{path: "statuses/destroy.json", method: methodPost}
  retweetEndpoint         = endpoint{path: "statuses/retweet.json", method: methodPost, group: publishLimitGroup}
  unretweetEndpoint       = endpoint{path: "statuses/unretweet.json", method: methodPost}
  likeEndpoint            = endpoint{path: "favorites/create.json", method: methodPost}
  unlikeEndpoint          = endpoint{path: "favorites/destroy.json", method: methodPost}
  updateProfileEndpoint   = endpoint{path: "account/update_profile.json", method: methodPost}
)

func (ep endpoint) limitKey() string {
  if ep.group != "" {
    return "group:" + string(ep.group)
  }
  return "singleton:" + ep.path
}

type twitterClient struct {
  client        *http.Client
  ses           *sessions
  protocol, url string
}

func newTwitterClient(timeout time.Duration, protocol, url string, assumeNextLimit bool) twitterClient {
  client := http.Client{
    Timeout: timeout,
  }
  return twitterClient{
    client:   &client,
    ses:      newSessions(assumeNextLimit),
    protocol: protocol,
    url:      url,
  }
}

func (tc twitterClient) standardRequest(ep endpoint, auth oauth.AuthPair, query, body oauth.Params, output interface{}) error {
  return tc.oauthRequest(ep, auth, query, body, func(resp *http.Response) error {
    return json.NewDecoder(resp.Body).Decode(output)
  })
}

func (tc twitterClient) oauthRequest(ep endpoint, auth oauth.AuthPair, query, body oauth.Params, handler func(resp *http.Response) error) error {
  oauthReq := oauth.NewRequest(ep.method.String(), tc.protocol, tc.url, ep.path, query, body)
  req, err := oauthReq.MakeRequest(auth)
  if err != nil {
    return err
  }
  return tc.request(req, ep, auth.Public.Token, handler)
}

func (tc twitterClient) request(req *http.Request, ep endpoint, token string, handler func(resp *http.Response) error) (err error) {
  resp, err := func() (*http.Response, error) {
    rl := tc.ses.get(token).getLimit(ep.limitKey())

    var (
      limitCurrent, limitNext *uint
      limitResets             *time.Time
      rateLimitHit            bool
    )

    if err := rl.use(); err != nil {
      return nil, err
    }

    defer func() {
      rl.finish(limitCurrent, limitNext, limitResets, rateLimitHit)
    }()

    resp, err := tc.client.Do(req)
    if err != nil {
      return nil, err //TODO: replace with custom error for connection failed
    }

    tooManyRequests := resp.StatusCode == http.StatusTooManyRequests

    var headerParseErr error

    if val, ok, err := parseLimitHeader(resp.Header.Get(headerRateLimitRemaining)); ok && err == nil {
      if !tooManyRequests {
        limitCurrent = new(uint)
        *limitCurrent = val
      }
    } else if err != nil {
      headerParseErr = err
    }

    if val, ok, err := parseLimitHeader(resp.Header.Get(headerRateLimit)); ok && err == nil {
      limitNext = new(uint)
      *limitNext = val
    } else if err != nil {
      headerParseErr = err
    }

    if val, ok, err := parseLimitResetsHeader(resp.Header.Get(headerRateLimitReset)); ok && err == nil {
      limitResets = new(time.Time)
      *limitResets = val
    } else if err != nil {
      headerParseErr = err
    }

    if tooManyRequests {
      rateLimitHit = true

      limitCurrent = new(uint)
      *limitCurrent = 0

      if limitResets != nil {
        return resp, newRateLimitError(*limitResets)
      } else {
        return resp, newRateLimitError(time.Time{})
      }
    }

    if headerParseErr != nil {
      return resp, newBadResponseError("Twitter responded with a rate limit header that could not be parsed")
    }

    return resp, nil
  }()

  if resp != nil {
    defer func() {
      if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
        err = closeErr
      }
    }()
  }

  if err != nil {
    return err
  }

  if 200 <= resp.StatusCode && resp.StatusCode < 300 {
    return handler(resp)
  } else {
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return err
    }
    msg := fmt.Sprintf("Twitter responded with %s: %s", resp.Status, string(body))
    if 400 <= resp.StatusCode && resp.StatusCode < 500 {
      return newBadRequestError(msg)
    } else {
      return newTwitterError(msg)
    }
  }
}

func parseLimitHeader(s string) (uint, bool, error) {
  if s == "" {
    return 0, false, nil
  }
  val, err := strconv.ParseUint(s, 10, bits.UintSize)
  if err != nil {
    return 0, false, err
  }
  return uint(val), true, nil
}

func parseLimitResetsHeader(s string) (time.Time, bool, error) {
  if s == "" {
    return time.Time{}, false, nil
  }
  unix, err := strconv.ParseInt(s, 10, 64)
  if err != nil {
    return time.Time{}, false, err
  }
  return time.Unix(unix, 0), true, nil
}
