package goldcrest

import (
  "fmt"
  "time"
)

type RateLimitError struct {
  resets time.Time
}

func (err RateLimitError) Error() string {
  return fmt.Sprintf("rate limit hit; retry at unix time %d", err.resets.Unix())
}

func (err RateLimitError) ResetsTime() time.Time {
  return err.resets
}

func (err RateLimitError) WaitDuration() time.Duration {
  return err.resets.Sub(err.resets)
}

type AmbiguousRateLimitError struct{}

func (err AmbiguousRateLimitError) Error() string {
  return "rate limit hit and reset time unknown"
}
