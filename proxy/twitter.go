package proxy

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/pantonshire/goldcrest/proxy/model"
  "github.com/pantonshire/goldcrest/proxy/oauth"
  "math/bits"
  "net/http"
  "path"
  "strconv"
  "strings"
  "time"
)

const (
  version  = "1.1"
  protocol = "https"
  domain   = "api.twitter.com"

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
)

func (ep endpoint) limitKey() string {
  if ep.group != "" {
    return "group:" + string(ep.group)
  }
  return "singleton:" + ep.path
}

type twitterClient struct {
  client *http.Client
  ses    *sessions
}

func newTwitterClient() twitterClient {
  client := http.Client{
    Timeout: time.Second * time.Duration(5 * time.Second), //TODO: put in config
  }
  return twitterClient{
    client: &client,
    ses:    newSessions(),
  }
}

func (tc twitterClient) getTweet(ctx context.Context, auth oauth.AuthPair, id uint64, opts tweetOptions) (model.Tweet, error) {
  query := opts.encode()
  query["id"] = strconv.FormatUint(id, 10)
  var tweet model.Tweet
  if err := tc.standardRequest(ctx, showTweetEndpoint, auth, query, nil, &tweet); err != nil {
    return model.Tweet{}, err
  }
  return tweet, nil
}

func (tc twitterClient) getTweets(ctx context.Context, auth oauth.AuthPair, ids []uint64, opts tweetOptions) (model.Timeline, error) {
  query := opts.encode()
  if len(ids) > 0 {
    idStrs := make([]string, len(ids))
    for i, id := range ids {
      idStrs[i] = strconv.FormatUint(id, 10)
    }
    query["id"] = strings.Join(idStrs, ",")
  }
  var tweets model.Timeline
  if err := tc.standardRequest(ctx, showTweetsEndpoint, auth, query, nil, &tweets); err != nil {
    return nil, err
  }
  return tweets, nil
}

func (tc twitterClient) standardRequest(ctx context.Context, ep endpoint, auth oauth.AuthPair, query, body map[string]string, output interface{}) error {
  return tc.oauthRequest(ctx, ep, auth, query, body, func(resp *http.Response) error {
    return json.NewDecoder(resp.Body).Decode(output)
  })
}

func (tc twitterClient) oauthRequest(ctx context.Context, ep endpoint, auth oauth.AuthPair, query, body map[string]string, handler func(resp *http.Response) error) error {
  oauthReq := oauth.NewRequest(ep.method.String(), protocol, domain, path.Join(version, ep.path), query, body)
  req, err := oauthReq.MakeRequest(ctx, auth.Secret, auth.Public)
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

    limitCurrent, limitNext, limitResets, err = rateLimitHeaders(resp.Header)
    if err != nil {
      return resp, err
    }

    if resp.StatusCode == http.StatusTooManyRequests {
      rateLimitHit = true
      if limitResets != nil {
        return resp, newRateLimitError(*limitResets)
      } else {
        return resp, newRateLimitError(time.Time{})
      }
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
  } else if 400 <= resp.StatusCode && resp.StatusCode < 500 {
    return newBadRequestError(fmt.Sprintf("Twitter responded with %s", resp.Status))
  } else {
    return newTwitterError(fmt.Sprintf("Twitter responded with %s", resp.Status))
  }
}

func rateLimitHeaders(header http.Header) (current, next *uint, resets *time.Time, err error) {
  if currentStr := header.Get(headerRateLimitRemaining); currentStr != "" {
    if currentVal, parseErr := strconv.ParseUint(currentStr, 10, bits.UintSize); parseErr != nil {
      err = newBadResponseHeaderError(headerRateLimitRemaining, currentStr)
    } else {
      current = new(uint)
      *current = uint(currentVal)
    }
  }
  if nextStr := header.Get(headerRateLimit); nextStr != "" {
    if nextVal, parseErr := strconv.ParseUint(nextStr, 10, bits.UintSize); parseErr != nil {
      err = newBadResponseHeaderError(headerRateLimit, nextStr)
    } else {
      next = new(uint)
      *next = uint(nextVal)
    }
  }
  if resetsStr := header.Get(headerRateLimitReset); resetsStr != "" {
    if resetsUnix, parseErr := strconv.ParseInt(resetsStr, 10, 64); parseErr != nil {
      err = newBadResponseHeaderError(headerRateLimitReset, resetsStr)
    } else {
      resets = new(time.Time)
      *resets = time.Unix(resetsUnix, 0)
    }
  }
  return current, next, resets, err
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
