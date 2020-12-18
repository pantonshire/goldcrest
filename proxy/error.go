package proxy

import (
  "fmt"
  pb "github.com/pantonshire/goldcrest/protocol"
  "google.golang.org/grpc/metadata"
  "strconv"
  "time"
)

type proxyError interface {
  error
  ser() (msg *pb.Error, meta metadata.MD)
}

type rateLimitError struct {
  retry time.Time
}

func newRateLimitError(retry time.Time) rateLimitError {
  return rateLimitError{retry: retry}
}

func (err rateLimitError) Error() string {
  return fmt.Sprintf("rate limit exceeded; resets at %s", err.retry.Format("15:04:05 MST"))
}

func (err rateLimitError) ser() (*pb.Error, metadata.MD) {
  var meta metadata.MD
  if !err.retry.IsZero() {
    meta = metadata.Pairs("retry", strconv.FormatInt(err.retry.Unix(), 10))
  }
  return &pb.Error{
    Code:    pb.Error_RATE_LIMIT,
    Message: err.Error(),
  }, meta
}

type twitterError struct {
  message string
}

func newTwitterError(message string) twitterError {
  return twitterError{message: message}
}

func (err twitterError) Error() string {
  return fmt.Sprintf("twitter connection error: %s", err.message)
}

func (err twitterError) ser() (*pb.Error, metadata.MD) {
  return &pb.Error{
    Code:    pb.Error_TWITTER_ERROR,
    Message: err.Error(),
  }, nil
}

type badRequestError struct {
  message string
}

func newBadRequestError(message string) badRequestError {
  return badRequestError{message: message}
}

func (err badRequestError) Error() string {
  return fmt.Sprintf("proxy sent bad request: %s", err.message)
}

func (err badRequestError) ser() (*pb.Error, metadata.MD) {
  return &pb.Error{
    Code:    pb.Error_BAD_REQUEST,
    Message: err.Error(),
  }, nil
}

type badResponseError struct {
  message string
}

func newBadResponseError(message string) badResponseError {
  return badResponseError{message: message}
}

func newBadResponseHeaderError(header, value string) badResponseError {
  return newBadResponseError(fmt.Sprintf("bad header value for %s: \"%s\"", header, value))
}

func (err badResponseError) Error() string {
  return fmt.Sprintf("twitter returned bad response: %s", err.message)
}

func (err badResponseError) ser() (*pb.Error, metadata.MD) {
  return &pb.Error{
    Code:    pb.Error_BAD_RESPONSE,
    Message: err.Error(),
  }, nil
}
