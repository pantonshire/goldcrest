package goldcrest

import (
  "fmt"
  "net/http"
)

type HttpError struct {
  Code   int
  Status string
}

func (e *HttpError) Error() string {
  return fmt.Sprintf("%d %s", e.Code, e.Status)
}

func NewHttpError(code int, status string) *HttpError {
  if IsStatusOK(code) {
    return nil
  }
  return &HttpError{Code: code, Status: status}
}

func HttpErrorFor(resp *http.Response) *HttpError {
  return NewHttpError(resp.StatusCode, resp.Status)
}

func IsStatusOK(code int) bool {
  return 200 <= code && code < 300
}
