package proxy

import (
  "sync"
  "time"
)

//TODO: put in config
var (
  stopOnLimitUnknown = true
  assumeNextLimit    = true
)

const shortWait = time.Millisecond * 50

type sessions struct {
  mx    sync.Mutex
  cache map[string]*session
}

type session struct {
  mx     sync.Mutex
  limits map[string]*rateLimit //Use endpoint.limitKey() as map key
}

//TODO: rl can get stuck on "stopped" if it never receives a response from Twitter, so remember to unblock
type rateLimit struct {
  mxData        sync.Mutex
  mxNext        sync.Mutex
  mxLow         sync.Mutex
  current       *uint
  next          *uint
  stopped       bool
  unknownUsages uint
  resets        time.Time
}

func newSession() *session {
  return &session{
    limits: make(map[string]*rateLimit),
  }
}

func newRateLimit() *rateLimit {
  return &rateLimit{}
}

func (ses *sessions) get(token string) *session {
  ses.mx.Lock()
  defer ses.mx.Unlock()
  if se, ok := ses.cache[token]; ok {
    return se
  }
  se := newSession()
  ses.cache[token] = se
  return se
}

func (se *session) getLimit(key string) *rateLimit {
  se.mx.Lock()
  defer se.mx.Unlock()
  if rl, ok := se.limits[key]; ok {
    return rl
  }
  rl := newRateLimit()
  se.limits[key] = rl
  return rl
}

func (rl *rateLimit) use() (bool, error) {
  //Lock with low priority
  rl.mxLow.Lock()
  defer rl.mxLow.Unlock()
  rl.mxNext.Lock()
  rl.mxData.Lock()
  defer rl.mxData.Unlock()
  rl.mxNext.Unlock()

  now := time.Now()
  reset := now.After(rl.resets)

  if reset && rl.next != nil {
    if rl.current == nil {
      rl.current = new(uint)
    }
    *rl.current = *rl.next
    rl.stopped = false
    if !assumeNextLimit {
      rl.next = nil
    }
  }

  if rl.current != nil && *rl.current > 0 {
    *rl.current--
    return false, nil
  }

  if rl.current == nil || rl.next == nil {
    if stopOnLimitUnknown {
      if !rl.stopped {
        rl.stopped = true
        rl.unknownUsages++
        return true, nil
      }
    } else {
      rl.unknownUsages++
      return false, nil
    }
  }

  var retry time.Time
  if reset {
    retry = now.Add(shortWait)
  } else {
    retry = rl.resets
  }
  return false, newRateLimitError(retry)
}

// forceSync should be used when a rate limit error occurred (a rate limit error occurring indicates that
// the proxy's rate limit tracker is wrong!), forcing the rateLimit to set all of its tracking fields.
func (rl *rateLimit) update(current, next uint, resets time.Time, forceSync bool) {
  //Lock with high priority
  rl.mxNext.Lock()
  rl.mxData.Lock()
  defer rl.mxData.Unlock()
  rl.mxNext.Unlock()

  if forceSync || rl.current == nil {
    if rl.current == nil {
      rl.current = new(uint)
    }
    *rl.current = current
  }

  if rl.next == nil {
    rl.next = new(uint)
  }

  *rl.next = next
  if resets.After(rl.resets) {
    rl.resets = resets
  }
}

func (rl *rateLimit) unblock() {
  //Lock with high priority
  rl.mxNext.Lock()
  rl.mxData.Lock()
  defer rl.mxData.Unlock()
  rl.mxNext.Unlock()

  rl.stopped = false
}
