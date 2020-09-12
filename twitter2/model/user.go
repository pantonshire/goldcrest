package model

import "time"

type User struct {
  ID              string    `json:"id"`
  DisplayName     string    `json:"name"`
  Username        string    `json:"username"`
  CreatedAt       time.Time `json:"created_at"`
  Description     string    `json:"description"`
  Location        string    `json:"location"`
  PinnedTweetID   string    `json:"pinned_tweet_id"`
  ProfileImageURL string    `json:"profile_image_url"`
  Protected       bool      `json:"protected"`
  URL             string    `json:"url"`
  Verified        bool      `json:"verified"`

  Entities struct {
    URL struct {
      URLs []URL `json:"urls"`
    } `json:"url"`
    Description struct {
      URLs     []URL `json:"urls"`
      Hashtags []Tag `json:"hashtags"`
      Mentions []Tag `json:"mentions"`
      Cashtags []Tag `json:"cashtags"`
    } `json:"description"`
  } `json:"entities"`

  PublicMetrics struct {
    FollowersCount int `json:"followers_count"`
    FollowingCount int `json:"following_count"`
    TweetCount     int `json:"tweet_count"`
    ListedCount    int `json:"listed_count"`
  } `json:"public_metrics"`
}
