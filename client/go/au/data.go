package au

import (
  "time"
)

type Tweet struct {
  ID                   uint64
  CreatedAt            time.Time
  Text                 string
  TextDisplayRange     Indices
  Truncated            bool
  Source               string
  User                 User
  RepliedTo            *ReplyData
  Quoted               *Tweet
  Retweeted            *Tweet
  Quotes               uint
  Replies              uint
  Retweets             uint
  Likes                uint
  CurrentUserLiked     bool
  CurrentUserRetweeted bool
  CurrentUserRetweetID *uint64
  Hashtags             []Symbol
  URLs                 []URL
  Mentions             []Mention
  Symbols              []Symbol
  Media                []Media
  Polls                []Poll
  PossiblySensitive    bool
  FilterLevel          string
  Lang                 string
  WithheldCopyright    bool
  WithheldCounties     []string
  WithheldScope        string
}

type ReplyData struct {
  TweetID    uint64
  UserID     uint64
  UserHandle string
}

type User struct {
  ID                  uint64
  Handle              string
  DisplayName         string
  CreatedAt           time.Time
  Bio                 string
  URL                 string
  Location            string
  Protected           bool
  Verified            bool
  FollowerCount       uint
  FollowingCount      uint
  ListedCount         uint
  FavouritesCount     uint
  StatusesCount       uint
  ProfileBanner       string
  ProfileImage        string
  DefaultProfile      bool
  DefaultProfileImage bool
  WithheldCountries   []string
  WithheldScope       string
  URLs                []URL
  BioURLs             []URL
}

type Indices struct {
  Start uint
  End   uint
}

type URL struct {
  Indices
  TwitterURL  string
  DisplayURL  string
  ExpandedURL string
}

type Symbol struct {
  Indices
  Text string
}

type Mention struct {
  Indices
  UserID          uint64
  UserHandle      string
  UserDisplayName string
}

type Media struct {
  URL
  ID            uint64
  Type          string
  MediaURL      string
  Alt           string
  SourceTweetID *uint64
  Thumb         MediaSize
  Small         MediaSize
  Medium        MediaSize
  Large         MediaSize
}

type MediaSize struct {
  Width  uint
  Height uint
  Resize string
}

type Poll struct {
  EndTime  time.Time
  Duration time.Duration
  Options  []PollOption
}

type PollOption struct {
  Position uint
  Text     string
}
