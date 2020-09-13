package goldcrest

import (
  "fmt"
  "net/http"
  "time"
)

type ClientConfig struct {
  RequestPoolSize uint `json:"pool"`
  TimeoutSeconds  uint `json:"timeoutSeconds"`
}

type client struct {
  http.Client
  pool chan int
}

type HttpError struct {
  code   int
  status string
}

func (e *HttpError) Error() string {
  return fmt.Sprintf("%d %s", e.code, e.status)
}

func NewHttpError(code int, status string) *HttpError {
  if IsStatusOK(code) {
    return nil
  }
  return &HttpError{code: code, status: status}
}

func HttpErrorFor(resp *http.Response) *HttpError {
  return NewHttpError(resp.StatusCode, resp.Status)
}

func IsStatusOK(code int) bool {
  return 200 <= code && code < 300
}

func newClient(conf ClientConfig) client {
  poolSize := int(conf.RequestPoolSize)
  t := client{
    Client: http.Client{
      Timeout: time.Second * time.Duration(conf.TimeoutSeconds),
    },
    pool: make(chan int, poolSize),
  }
  for i := 0; i < poolSize; i++ {
    t.pool <- i
  }
  return t
}

func (t client) close() {
  close(t.pool)
}

func (t client) Do(r *http.Request) (*http.Response, error) {
  rid := <-t.pool
  defer func() { t.pool <- rid }()
  return t.Client.Do(r)
}
