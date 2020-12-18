package proxy

import (
  "sync"
  "time"
)

const (
  stuckResetTime = time.Minute * 20
)

type sessions struct {
  mx    sync.Mutex
  cache map[string]*session
}

type session struct {
  mx     sync.Mutex
  limits map[string]*rateLimit //Use endpoint.limitKey() as map key
}

type rateLimit struct {
  mxData      sync.Mutex
  mxNext      sync.Mutex
  mxLow       sync.Mutex
  resolving   bool
  mxResolving sync.Mutex
  current     *uint
  next        *uint
  resets      time.Time
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

func (rl *rateLimit) lockLow() {
  rl.mxLow.Lock()
  rl.mxNext.Lock()
  rl.mxData.Lock()
  rl.mxNext.Unlock()
}

func (rl *rateLimit) unlockLow() {
  rl.mxData.Unlock()
  rl.mxLow.Unlock()
}

func (rl *rateLimit) lockHigh() {
  rl.mxNext.Lock()
  rl.mxData.Lock()
  rl.mxNext.Unlock()
}

func (rl *rateLimit) unlockHigh() {
  rl.mxData.Unlock()
}

func (rl *rateLimit) use() error {
  rl.lockLow()
  defer rl.unlockLow()

  for rl.resolving {
    rl.unlockLow()
    rl.mxResolving.Lock()
    rl.mxResolving.Unlock()
    rl.lockLow()
  }

  now := time.Now()
  resetsKnown := !rl.resets.IsZero()

  if !resetsKnown && rl.current != nil && *rl.current == 0 {
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
      rl.next = nil
    }
    rl.resets = time.Time{}
  }

  if rl.current == nil {
    rl.mxResolving.Lock()
    rl.resolving = true
  } else if *rl.current > 0 {
    *rl.current--
    return nil
  }

  return newRateLimitError(rl.resets)
}

func (rl *rateLimit) finish(current, next *uint, resets *time.Time, forceSync bool) {
  rl.lockHigh()
  defer rl.unlockHigh()

  if rl.resolving {
    rl.mxResolving.Unlock()
  }

  if current != nil && forceSync || rl.current == nil {
    if rl.current == nil {
      rl.current = new(uint)
    }
    *rl.current = *current
  }

  if next != nil {
    if rl.next == nil {
      rl.next = new(uint)
    }
    *rl.next = *next
  }

  if resets != nil && resets.After(rl.resets) {
    rl.resets = *resets
  }
}
