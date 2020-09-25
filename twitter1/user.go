package twitter1

import (
  "context"
  "fmt"
  "net/http"
  "strconv"
  "sync"
  "time"
)

type limitGroup string

const (
  xRateLimit          = "X-Rate-Limit-Limit"
  xRateLimitRemaining = "X-Rate-Limit-Remaining"
  xRateLimitReset     = "X-Rate-Limit-Reset"

  limitNone            limitGroup = ""
  limitStatusUpdate    limitGroup = "statuses/update"
  limitStatusShow      limitGroup = "statuses/show"
  limitHomeTimeline    limitGroup = "statuses/home_timeline"
  limitMentionTimeline limitGroup = "statuses/mentions_timeline"
  limitUserTimeline    limitGroup = "statuses/user_timeline"
  limitUpdateProfile   limitGroup = "statuses/update_profile"
)

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
  current, next, resets, err := getRateLimitHeaders(resp.Header)
  if current != nil {
    rl.current = *current
  }
  if next != nil {
    rl.next = *next
  }
  if resets != nil {
    rl.resets = *resets
  }
  if err != nil {
    return resp, err
  }
  return resp, nil
}

func getRateLimitHeaders(header http.Header) (current, next *int, resets *time.Time, err error) {
  if currentStr := header.Get(xRateLimitRemaining); currentStr != "" {
    currentVal, parseErr := strconv.Atoi(currentStr)
    if parseErr != nil {
      err = fmt.Errorf("invalid rate limit header for %s: \"%s\"", xRateLimitRemaining, currentStr)
    }
    current = &currentVal
  }
  if nextStr := header.Get(xRateLimit); nextStr != "" {
    nextVal, parseErr := strconv.Atoi(nextStr)
    if parseErr != nil {
      err = fmt.Errorf("invalid rate limit header for %s: \"%s\"", xRateLimit, nextStr)
    }
    next = &nextVal
  }
  if resetsStr := header.Get(xRateLimitReset); resetsStr != "" {
    resetsUnix, parseErr := strconv.ParseInt(resetsStr, 10, 64)
    if parseErr != nil {
      err = fmt.Errorf("invalid rate limit header for %s: \"%s\"", xRateLimitReset, resetsStr)
    }
    resetsVal := time.Unix(resetsUnix, 0)
    resets = &resetsVal
  }
  return current, next, resets, err
}
