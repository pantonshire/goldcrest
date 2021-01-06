package proxy

import (
  "sync"
  "time"
)

const (
  stuckResetTime = time.Minute * 20
)

type sessions struct {
  mx              sync.Mutex
  cache           map[string]*session
  assumeNextLimit bool
}

type session struct {
  mx              sync.Mutex
  limits          map[string]*rateLimit
  assumeNextLimit bool
}

type rateLimit struct {
  mxData     sync.Mutex
  mxNext     sync.Mutex
  mxLow      sync.Mutex
  //boolean flag is safe from instruction reordering because it is always protected by mxData
  resolving  bool
  resolved   chan struct{}
  current    *uint
  next       *uint
  resets     time.Time
  assumeNext bool
}

func newSessions(assumeNextLimit bool) *sessions {
  return &sessions{
    cache:           make(map[string]*session),
    assumeNextLimit: assumeNextLimit,
  }
}

func newSession(assumeNextLimit bool) *session {
  return &session{
    limits:          make(map[string]*rateLimit),
    assumeNextLimit: assumeNextLimit,
  }
}

func newRateLimit(assumeNext bool) *rateLimit {
  resolved := make(chan struct{}, 1)
  resolved <- struct{}{}
  return &rateLimit{
    resolved:   resolved,
    assumeNext: assumeNext,
  }
}

func (ses *sessions) get(token string) *session {
  ses.mx.Lock()
  defer ses.mx.Unlock()
  if se, ok := ses.cache[token]; ok {
    return se
  }
  se := newSession(ses.assumeNextLimit)
  ses.cache[token] = se
  return se
}

func (se *session) getLimit(key string) *rateLimit {
  se.mx.Lock()
  defer se.mx.Unlock()
  if rl, ok := se.limits[key]; ok {
    return rl
  }
  rl := newRateLimit(se.assumeNextLimit)
  se.limits[key] = rl
  return rl
}

func (rl *rateLimit) lockLow() {
  log.Debug("Wait for low lock")
  rl.mxLow.Lock()
  rl.mxNext.Lock()
  rl.mxData.Lock()
  rl.mxNext.Unlock()
  log.Debug("Acquire low lock")
}

func (rl *rateLimit) unlockLow() {
  rl.mxData.Unlock()
  rl.mxLow.Unlock()
  log.Debug("Release low lock")
}

func (rl *rateLimit) lockHigh() {
  log.Debug("Wait for high lock")
  rl.mxNext.Lock()
  rl.mxData.Lock()
  rl.mxNext.Unlock()
  log.Debug("Acquire high lock")
}

func (rl *rateLimit) unlockHigh() {
  rl.mxData.Unlock()
  log.Debug("Release high lock")
}

func (rl *rateLimit) use() error {
  rl.lockLow()
  defer rl.unlockLow()

  for rl.resolving {
    log.Debug("Resolving! Must wait")
    rl.unlockLow()
    log.Debug("Wait for resolved message")
    <-rl.resolved
    log.Debug("Received resolved message")
    rl.resolved <- struct{}{}
    log.Debug("Return resolved message")
    rl.lockLow()
  }

  now := time.Now()
  resetsKnown := !rl.resets.IsZero()

  if !resetsKnown && rl.current != nil && *rl.current == 0 {
    log.Info("Escape from rate limit stuck condition")
    rl.resets = now.Add(stuckResetTime) //Could be stuck forever otherwise!
  } else if resetsKnown && now.After(rl.resets) {
    if rl.next == nil {
      if rl.current != nil && *rl.current == 0 {
        rl.current = nil
      }
    } else {
      if rl.current == nil {
        rl.current = new(uint)
      }
      *rl.current = *rl.next
      if !rl.assumeNext {
        rl.next = nil
      }
    }
    rl.resets = time.Time{}
  }

  if rl.current == nil {
    log.Debug("Start resolving, take from resolved message channel")
    <-rl.resolved
    log.Debug("Received resolved message")
    rl.resolving = true
    return nil
  } else if *rl.current > 0 {
    log.WithField("old", *rl.current).WithField("new", *rl.current-1).Info("Update limit")
    *rl.current--
    return nil
  } else {
    log.Info("Rate limit error")
    return newRateLimitError(rl.resets)
  }
}

func (rl *rateLimit) finish(current, next *uint, resets *time.Time, forceSync bool) {
  rl.lockHigh()
  defer rl.unlockHigh()

  if rl.resolving {
    log.Debug("Finished resolving, send resolved message")
    rl.resolving = false
    rl.resolved <- struct{}{}
  }

  if current != nil && (forceSync || rl.current == nil) {
    if rl.current == nil {
      rl.current = new(uint)
    }
    *rl.current = *current
    log.WithField("value", *current).Info("Get new current limit")
  }

  if next != nil {
    if rl.next == nil {
      rl.next = new(uint)
    }
    *rl.next = *next
    log.WithField("value", *next).Info("Get new next limit")
  }

  if resets != nil && resets.After(rl.resets) {
    rl.resets = *resets
    log.WithField("value", *resets).Info("Get new resets time")
  }
}
