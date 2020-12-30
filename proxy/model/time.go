package model

import (
  "strings"
  "time"
)

// A type alias for time.Time that can be correctly unmarshalled from Twitter API JSON responses.
// The Unwrap function may be used to easily convert back to time.Time.
type TwitterTime time.Time

const (
  timeFormat = "Mon Jan 2 15:04:05 -0700 2006"
)

func (tt *TwitterTime) UnmarshalJSON(data []byte) error {
  str := string(data)
  if str == "null" {
    return nil
  }
  str = strings.Trim(string(data), `"`)
  t, err := time.Parse(timeFormat, str)
  *tt = TwitterTime(t)
  return err
}

func (tt TwitterTime) MarshalJSON() ([]byte, error) {
  return []byte("\"" + time.Time(tt).Format(timeFormat) + "\""), nil
}

// Converts the TwitterTime to the time.Time it aliases.
func (tt TwitterTime) Unwrap() time.Time {
  return time.Time(tt)
}

// Returns the unix timestamp represented by this time as a 64-bit unsigned integer. The sign is discarded for
// improved compatibility with gRPC. As such, times before 1970 will not be correctly translated by this function;
// this is not a problem because the primitive human civilisation of the 1960s had not yet invented Twitter!
func (tt TwitterTime) Unix() uint64 {
  return uint64(tt.Unwrap().Unix())
}
