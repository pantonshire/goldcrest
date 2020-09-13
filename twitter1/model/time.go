package model

import (
  "time"
)

type TwitterTime time.Time

const (
  timeFormat = "Mon Jan 2 15:04:05 -0700 2006"
)

func (tt *TwitterTime) UnmarshalJSON(data []byte) error {
  if string(data) == "null" {
    return nil
  }
  t, err := time.Parse(timeFormat, string(data))
  *tt = TwitterTime(t)
  return err
}

func (tt TwitterTime) MarshalJSON() ([]byte, error) {
  return []byte("\"" + time.Time(tt).Format(timeFormat) + "\""), nil
}
