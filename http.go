package goldcrest

import (
  "fmt"
  "net/http"
  "strings"
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
  var status string
  if splitStatus := strings.SplitN(resp.Status, " ", 2); len(splitStatus) > 1 {
    status = splitStatus[1]
  }
  return NewHttpError(resp.StatusCode, status)
}

func IsStatusOK(code int) bool {
  return 200 <= code && code < 300
}
