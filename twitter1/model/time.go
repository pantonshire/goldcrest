package model

import (
  "strings"
  "time"
)

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
