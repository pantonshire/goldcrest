package au

import (
  pb "github.com/pantonshire/goldcrest/protocol"
  "time"
)

func desTimeline(msg *pb.Tweets) []Tweet {
  if msg == nil {
    return nil
  }
  tweets := make([]Tweet, len(msg.Tweets))
  for i, tweetMsg := range msg.Tweets {
    tweets[i] = desTweet(tweetMsg)
  }
  return tweets
}

func desTweet(msg *pb.Tweet) Tweet {
  if msg == nil {
    return Tweet{}
  }
  tweet := Tweet{
    ID:                   msg.Id,
    CreatedAt:            time.Unix(int64(msg.CreatedAt), 0),
    Text:                 msg.Text,
    TextDisplayRange:     desIndices(msg.TextDisplayRange),
    Truncated:            msg.Truncated,
    Source:               msg.Source,
    User:                 desUser(msg.User),
    Quotes:               uint(msg.QuoteCount),
    Replies:              uint(msg.ReplyCount),
    Retweets:             uint(msg.RetweetCount),
    Likes:                uint(msg.FavoriteCount),
    CurrentUserLiked:     msg.Favorited,
    CurrentUserRetweeted: msg.Retweeted,
    Hashtags:             desSymbols(msg.Hashtags),
    URLs:                 desURLs(msg.Urls),
    Mentions:             desMentions(msg.Mentions),
    Symbols:              desSymbols(msg.Symbols),
    Media:                desMediaItems(msg.Media),
    Polls:                desPolls(msg.Polls),
    PossiblySensitive:    msg.PossiblySensitive,
    FilterLevel:          msg.FilterLevel,
    Lang:                 msg.Lang,
    WithheldCopyright:    msg.WithheldCopyright,
    WithheldCounties:     msg.WithheldCountries,
    WithheldScope:        msg.WithheldScope,
  }
  if reply, ok := msg.Reply.(*pb.Tweet_RepliedTweet); ok && reply.RepliedTweet != nil {
    tweet.RepliedTo = &ReplyData{
      TweetID:    reply.RepliedTweet.ReplyToTweetId,
      UserID:     reply.RepliedTweet.ReplyToUserId,
      UserHandle: reply.RepliedTweet.ReplyToUserHandle,
    }
  }
  if quote, ok := msg.Quote.(*pb.Tweet_QuotedTweet); ok {
    tweet.Quoted = new(Tweet)
    *tweet.Quoted = desTweet(quote.QuotedTweet)
  }
  if retweet, ok := msg.Retweet.(*pb.Tweet_RetweetedTweet); ok {
    tweet.Retweeted = new(Tweet)
    *tweet.Retweeted = desTweet(retweet.RetweetedTweet)
  }
  if currentUserRetweet, ok := msg.CurrentUserRetweet.(*pb.Tweet_CurrentUserRetweetId); ok {
    tweet.CurrentUserRetweetID = new(uint64)
    *tweet.CurrentUserRetweetID = currentUserRetweet.CurrentUserRetweetId
  }
  return tweet
}

func desUser(msg *pb.User) User {
  if msg == nil {
    return User{}
  }
  return User{
    ID:                  msg.Id,
    Handle:              msg.Handle,
    DisplayName:         msg.DisplayName,
    CreatedAt:           time.Unix(int64(msg.CreatedAt), 0),
    Bio:                 msg.Bio,
    URL:                 msg.Url,
    Location:            msg.Location,
    Protected:           msg.Protected,
    Verified:            msg.Verified,
    FollowerCount:       uint(msg.FollowerCount),
    FollowingCount:      uint(msg.FollowingCount),
    ListedCount:         uint(msg.ListedCount),
    FavouritesCount:     uint(msg.FavoritesCount),
    StatusesCount:       uint(msg.StatusesCount),
    ProfileBanner:       msg.ProfileBanner,
    ProfileImage:        msg.ProfileImage,
    DefaultProfile:      msg.DefaultProfile,
    DefaultProfileImage: msg.DefaultProfileImage,
    WithheldCountries:   msg.WithheldCountries,
    WithheldScope:       msg.WithheldScope,
    URLs:                desURLs(msg.UrlUrls),
    BioURLs:             desURLs(msg.BioUrls),
  }
}

func desIndices(msg *pb.Indices) Indices {
  if msg == nil {
    return Indices{}
  }
  return Indices{
    Start: uint(msg.Start),
    End:   uint(msg.End),
  }
}

func desURLs(msgs []*pb.URL) []URL {
  urls := make([]URL, len(msgs))
  for i, msg := range msgs {
    urls[i] = desURL(msg)
  }
  return urls
}

func desURL(msg *pb.URL) URL {
  if msg == nil {
    return URL{}
  }
  return URL{
    Indices:     desIndices(msg.Indices),
    TwitterURL:  msg.TwitterUrl,
    DisplayURL:  msg.DisplayUrl,
    ExpandedURL: msg.ExpandedUrl,
  }
}

func desSymbols(msgs []*pb.Symbol) []Symbol {
  symbols := make([]Symbol, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      symbols[i] = Symbol{
        Indices: desIndices(msg.Indices),
        Text:    msg.Text,
      }
    }
  }
  return symbols
}

func desMentions(msgs []*pb.Mention) []Mention {
  mentions := make([]Mention, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      mentions[i] = Mention{
        Indices:         desIndices(msg.Indices),
        UserID:          msg.UserId,
        UserHandle:      msg.Handle,
        UserDisplayName: msg.DisplayName,
      }
    }
  }
  return mentions
}

func desMediaItems(msgs []*pb.Media) []Media {
  media := make([]Media, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      media[i] = Media{
        URL:           desURL(msg.Url),
        ID:            msg.Id,
        Type:          msg.Type,
        MediaURL:      msg.MediaUrl,
        Alt:           msg.Alt,
        SourceTweetID: desMediaSource(msg),
        Thumb:         desMediaSize(msg.Thumb),
        Small:         desMediaSize(msg.Small),
        Medium:        desMediaSize(msg.Medium),
        Large:         desMediaSize(msg.Large),
      }
    }
  }
  return media
}

func desMediaSource(msg *pb.Media) *uint64 {
  if source, ok := msg.Source.(*pb.Media_SourceTweetId); ok {
    id := new(uint64)
    *id = source.SourceTweetId
    return id
  }
  return nil
}

func desMediaSize(msg *pb.Media_Size) MediaSize {
  if msg == nil {
    return MediaSize{}
  }
  return MediaSize{
    Width:  uint(msg.Width),
    Height: uint(msg.Height),
    Resize: msg.Resize,
  }
}

func desPolls(msgs []*pb.Poll) []Poll {
  polls := make([]Poll, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      polls[i] = Poll{
        EndTime:  time.Unix(int64(msg.EndTime), 0),
        Duration: time.Minute * time.Duration(msg.DurationMinutes),
        Options:  make([]PollOption, len(msg.Options)),
      }
      for j, optMsg := range msg.Options {
        if optMsg != nil {
          polls[i].Options[j] = PollOption{
            Position: uint(optMsg.Position),
            Text:     optMsg.Text,
          }
        }
      }
    }
  }
  return polls
}
