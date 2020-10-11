package twitter1

import (
  "strings"
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

func (tweet Tweet) Unfold() Tweet {
  if tweet.Retweeted != nil {
    return *tweet.Retweeted
  }
  return tweet
}

func (tweet Tweet) TextOnly() string {
  var removeIndices []Indices
  if !tweet.TextDisplayRange.IsZero() {
    removeIndices = append(removeIndices, tweet.TextDisplayRange.Invert(uint(len(tweet.Text)))...)
  }
  for _, hashtag := range tweet.Hashtags {
    if !hashtag.Indices.IsZero() {
      removeIndices = append(removeIndices, hashtag.Indices)
    }
  }
  for _, url := range tweet.URLs {
    if !url.Indices.IsZero() {
      removeIndices = append(removeIndices, url.Indices)
    }
  }
  for _, mention := range tweet.Mentions {
    if !mention.Indices.IsZero() {
      removeIndices = append(removeIndices, mention.Indices)
    }
  }
  for _, symbol := range tweet.Symbols {
    if !symbol.Indices.IsZero() {
      removeIndices = append(removeIndices, symbol.Indices)
    }
  }
  for _, media := range tweet.Media {
    if !media.Indices.IsZero() {
      removeIndices = append(removeIndices, media.Indices)
    }
  }
  return removeFromString(tweet.Text, removeIndices...)
}

func (tweet Tweet) ReplyText(text string) string {
  return tweet.User.AtHandle() + " " + text
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

func (user User) AtHandle() string {
  return "@" + user.Handle
}

type Indices struct {
  Start uint
  End   uint
}

func (indices Indices) IsZero() bool {
  return indices.Start <= indices.End
}

func (indices Indices) Invert(l uint) []Indices {
  var inverted []Indices
  if indices.Start > 0 {
    inverted = append(inverted, Indices{Start: 0, End: indices.Start})
  }
  if indices.End < l {
    inverted = append(inverted, Indices{Start: indices.End, End: l})
  }
  return inverted
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

func (mention Mention) AtHandle() string {
  return "@" + mention.UserHandle
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

func removeFromString(s string, cutAt ...Indices) string {
  n := uint(len(s))
  ignorePos := make([]bool, n)
  for _, indices := range cutAt {
    for i := indices.Start; i < indices.End && i < n; i++ {
      ignorePos[i] = true
    }
  }
  var buf strings.Builder
  for i, r := range s {
    if !ignorePos[i] {
      buf.WriteRune(r)
    }
  }
  return buf.String()
}
