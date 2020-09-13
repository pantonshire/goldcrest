package goldcrest

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
  "path"
  "sync"
)

type TwitterConfig struct {
  Client       ClientConfig `json:"client"`
  LimitPadding uint         `json:"limitPadding"`
}

type Twitter struct {
  c            client
  limitPadding uint
  bearerLock   sync.Mutex
  bearers      map[string]*twitterBearer
}

type twitterBearer struct {
  token     string
  limitLock sync.Mutex
  limit     *int
  nextLimit *int
}

const (
  twitterBaseURL = "https://api.twitter.com"
  twitterAPIv2   = "2"

  xRateLimit          = "X-Rate-Limit-Limit"
  xRateLimitRemaining = "X-Rate-Limit-Remaining"
  xRateLimitReset     = "X-Rate-Limit-Reset"
)

//TODO: rate limiting
func (t *Twitter) req(bearer *twitterBearer, method, reqPath, version string, body, response interface{}) (retry bool, err error) {
  b, err := json.Marshal(body)
  if err != nil {
    return false, err
  }
  req, err := http.NewRequest(method, path.Join(twitterBaseURL, version, reqPath), bytes.NewBuffer(b))
  if err != nil {
    return false, err
  }
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer.token))
  resp, err := t.c.Do(req)
  if err != nil {
    return false, err
  }
  defer func() {
    if closeErr := resp.Body.Close(); closeErr != nil {
      err = closeErr
    }
  }()
  httpErr := HttpErrorFor(resp)
  if err == nil {
    if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
      return false, err
    }
    return false, nil
  } else if httpErr.code == http.StatusTooManyRequests {
    return true, nil
  }
  return false, HttpErrorFor(resp)
}
