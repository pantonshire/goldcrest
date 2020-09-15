package twitter1

import (
  "sync"
  "time"
)

type rateLimit struct {
  lock          sync.Mutex
  current, next int
  resets        time.Time
}
