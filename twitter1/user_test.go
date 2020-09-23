package twitter1

import (
  "context"
  "fmt"
  "net/http"
  "strconv"
  "testing"
  "time"
)

func TestRateLimit(t *testing.T) {
  t.Run("RateLimit_1", func(t *testing.T) {
    oldResets := time.Now().Add(time.Second)
    newResets := oldResets.Add(time.Minute * 15)
    rl := rateLimit{current: 1, resets: oldResets}
    _, err := rl.do(context.Background(), func() (*http.Response, error) {
      return &http.Response{
        Header: map[string][]string{
          "X-Rate-Limit-Limit":     {strconv.Itoa(100)},
          "X-Rate-Limit-Remaining": {strconv.Itoa(0)},
          "X-Rate-Limit-Reset":     {fmt.Sprint(newResets.Unix())},
        },
      }, nil
    })
    if err != nil {
      t.Error(err)
    }
    if time.Now().After(oldResets) {
      t.Error("rate limit waited for reset despite having non-zero current value")
    }
    if rl.current != 0 {
      t.Errorf("rl.current is \"%d\", expected \"%d\"", rl.current, 0)
    }
    if rl.next != 100 {
      t.Errorf("rl.next is \"%d\", expected \"%d\"", rl.next, 100)
    }
    if rl.resets.Unix() != newResets.Unix() {
      t.Errorf("rl.resets is \"%d\", expected \"%d\"", rl.resets.Unix(), newResets.Unix())
    }
  })

  t.Run("RateLimit_2", func(t *testing.T) {
    resets := time.Now().Add(time.Millisecond * 50)
    rl := rateLimit{current: 0, resets: resets}
    _, err := rl.do(context.Background(), func() (*http.Response, error) {
      return &http.Response{
        Header: map[string][]string{
          "X-Rate-Limit-Limit":     {strconv.Itoa(100)},
          "X-Rate-Limit-Remaining": {strconv.Itoa(100)},
          "X-Rate-Limit-Reset":     {fmt.Sprint(resets.Add(time.Minute * 15).Unix())},
        },
      }, nil
    })
    if err != nil {
      t.Error(err)
    }
    if time.Now().Before(resets) {
      t.Error("rate limit did not wait for reset")
    }
    if rl.current != 100 {
      t.Errorf("rl.current is \"%d\", expected \"%d\"", rl.current, 100)
    }
  })
}
