package twitter1

import (
  pb "github.com/pantonshire/goldcrest/proto"
  "math"
  "time"
)

func encodeAuthPairMessage(auth AuthPair) *pb.Authentication {
  return &pb.Authentication{
    ConsumerKey: auth.Public.Key,
    AccessToken: auth.Public.Token,
    SecretKey:   auth.Secret.Key,
    SecretToken: auth.Secret.Token,
  }
}

func decodeAuthPairMessage(msg *pb.Authentication) AuthPair {
  if msg == nil {
    return AuthPair{}
  }
  return AuthPair{
    Auth{Key: msg.SecretKey, Token: msg.SecretToken},
    Auth{Key: msg.ConsumerKey, Token: msg.AccessToken},
  }
}

func encodeTweetOptionsMessage(twopts TweetOptions) *pb.TweetOptions {
  return &pb.TweetOptions{
    TrimUser:          twopts.TrimUser,
    IncludeMyRetweet:  twopts.IncludeMyRetweet,
    IncludeEntities:   twopts.IncludeEntities,
    IncludeExtAltText: twopts.IncludeExtAltText,
    IncludeCardUri:    twopts.TrimUser,
    Mode:              encodeTweetModeMessage(twopts.Mode),
  }
}

func decodeTweetOptionsMessage(msg *pb.TweetOptions) TweetOptions {
  if msg == nil {
    return TweetOptions{}
  }
  return TweetOptions{
    TrimUser:          msg.TrimUser,
    IncludeMyRetweet:  msg.IncludeMyRetweet,
    IncludeEntities:   msg.IncludeEntities,
    IncludeExtAltText: msg.IncludeExtAltText,
    IncludeCardURI:    msg.TrimUser,
    Mode:              decodeTweetModeMessage(msg.Mode),
  }
}

func encodeTimelineOptionsMessage(tlOpts TimelineOptions) *pb.TimelineOptions {
  msg := pb.TimelineOptions{}
  if tlOpts.Count != nil {
    msg.Count = uint32(*tlOpts.Count)
  }
  if tlOpts.MinID != nil {
    msg.MinId = *tlOpts.MinID
  }
  if tlOpts.MaxID != nil {
    msg.MaxId = *tlOpts.MaxID
  }
  return &msg
}

func decodeTimelineOptionsMessage(msg *pb.TimelineOptions) TimelineOptions {
  tlOpts := TimelineOptions{}
  if msg.Count != 0 {
    count := uint(msg.Count)
    tlOpts.Count = &count
  }
  if msg.MinId != 0 {
    minID := msg.MinId
    tlOpts.MinID = &minID
  }
  if msg.MaxId != 0 {
    maxID := msg.MaxId
    tlOpts.MaxID = &maxID
  }
  return tlOpts
}

func encodeTweetModeMessage(mode TweetMode) pb.TweetOptions_Mode {
  if mode == ExtendedMode {
    return pb.TweetOptions_EXTENDED
  }
  return pb.TweetOptions_COMPAT
}

func decodeTweetModeMessage(msg pb.TweetOptions_Mode) TweetMode {
  if msg == pb.TweetOptions_EXTENDED {
    return ExtendedMode
  }
  return CompatibilityMode
}

func decodeTimelineMessage(msg *pb.Timeline) []Tweet {
  if msg == nil {
    return nil
  }
  tweets := make([]Tweet, len(msg.Tweets))
  for i, tweetMsg := range msg.Tweets {
    tweets[i] = decodeTweetMessage(tweetMsg)
  }
  return tweets
}

func decodeTweetMessage(msg *pb.Tweet) Tweet {
  if msg == nil {
    return Tweet{}
  }
  tweet := Tweet{
    ID:                   msg.Id,
    CreatedAt:            time.Unix(msg.CreatedAt, 0),
    Text:                 msg.Text,
    TextDisplayRange:     decodeIndicesMessage(msg.TextDisplayRange),
    Truncated:            msg.Truncated,
    Source:               msg.Source,
    User:                 decodeUserMessage(msg.User),
    Quotes:               uint(msg.QuoteCount),
    Replies:              uint(msg.ReplyCount),
    Retweets:             uint(msg.RetweetCount),
    Likes:                uint(msg.FavoriteCount),
    CurrentUserLiked:     msg.Favorited,
    CurrentUserRetweeted: msg.Retweeted,
    Hashtags:             decodeSymbolMessages(msg.Hashtags),
    URLs:                 decodeURLMessages(msg.Urls),
    Mentions:             decodeMentionMessages(msg.Mentions),
    Symbols:              decodeSymbolMessages(msg.Symbols),
    Media:                decodeMediaMessages(msg.Media),
    Polls:                decodePolls(msg.Polls),
    PossiblySensitive:    msg.PossiblySensitive,
    FilterLevel:          msg.FilterLevel,
    Lang:                 msg.Lang,
    WithheldCopyright:    msg.WithheldCopyright,
    WithheldCounties:     msg.WithheldCountries,
    WithheldScope:        msg.WithheldScope,
  }
  if reply, ok := msg.Reply.(*pb.Tweet_RepliedTweet); ok && reply != nil {
    if reply.RepliedTweet != nil {
      tweet.RepliedTo = &ReplyData{
        TweetID:    reply.RepliedTweet.ReplyToTweetId,
        UserID:     reply.RepliedTweet.ReplyToUserId,
        UserHandle: reply.RepliedTweet.ReplyToUserHandle,
      }
    }
  }
  if quote, ok := msg.Quote.(*pb.Tweet_QuotedTweet); ok && quote != nil {
    decodedQuote := decodeTweetMessage(quote.QuotedTweet)
    tweet.Quoted = &decodedQuote
  }
  if retweet, ok := msg.Retweet.(*pb.Tweet_RetweetedTweet); ok && retweet != nil {
    decodedRetweet := decodeTweetMessage(retweet.RetweetedTweet)
    tweet.Retweeted = &decodedRetweet
  }
  if msg.CurrentUserRetweetId != 0 {
    retweetID := msg.CurrentUserRetweetId
    tweet.CurrentUserRetweetID = &retweetID
  }
  return tweet
}

func decodeUserMessage(msg *pb.User) User {
  if msg == nil {
    return User{}
  }
  return User{
    ID:                  msg.Id,
    Handle:              msg.Handle,
    DisplayName:         msg.DisplayName,
    CreatedAt:           time.Unix(msg.CreatedAt, 0),
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
    URLs:                decodeURLMessages(msg.UrlUrls),
    BioURLs:             decodeURLMessages(msg.BioUrls),
  }
}

func encodeIndicesMessage(indices Indices) *pb.Indices {
  return &pb.Indices{Start: uint32(indices.Start), End: uint32(indices.End)}
}

func decodeIndicesMessage(msg *pb.Indices) Indices {
  if msg == nil {
    return Indices{}
  }
  return Indices{Start: uint(msg.Start), End: uint(msg.End)}
}

func encodeURLMessages(urls []URL) []*pb.URL {
  msgs := make([]*pb.URL, len(urls))
  for i, url := range urls {
    msgs[i] = encodeURLMessage(url)
  }
  return msgs
}

func encodeURLMessage(url URL) *pb.URL {
  return &pb.URL{
    Indices:     encodeIndicesMessage(url.Indices),
    TwitterUrl:  url.TwitterURL,
    DisplayUrl:  url.DisplayURL,
    ExpandedUrl: url.ExpandedURL,
  }
}

func decodeURLMessages(msgs []*pb.URL) []URL {
  urls := make([]URL, len(msgs))
  for i, msg := range msgs {
    urls[i] = decodeURLMessage(msg)
  }
  return urls
}

func decodeURLMessage(msg *pb.URL) URL {
  if msg == nil {
    return URL{}
  }
  return URL{
    Indices:     decodeIndicesMessage(msg.Indices),
    TwitterURL:  msg.TwitterUrl,
    DisplayURL:  msg.DisplayUrl,
    ExpandedURL: msg.ExpandedUrl,
  }
}

func encodeSymbolMessages(symbols []Symbol) []*pb.Symbol {
  msgs := make([]*pb.Symbol, len(symbols))
  for i, symbol := range symbols {
    msgs[i] = encodeSymbolMessage(symbol)
  }
  return msgs
}

func encodeSymbolMessage(symbol Symbol) *pb.Symbol {
  return &pb.Symbol{
    Indices: encodeIndicesMessage(symbol.Indices),
    Text:    symbol.Text,
  }
}

func decodeSymbolMessages(msgs []*pb.Symbol) []Symbol {
  symbols := make([]Symbol, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      symbols[i] = Symbol{
        Indices: decodeIndicesMessage(msg.Indices),
        Text:    msg.Text,
      }
    }
  }
  return symbols
}

func encodeMentionMessages(mentions []Mention) []*pb.Mention {
  msgs := make([]*pb.Mention, len(mentions))
  for i, mention := range mentions {
    msgs[i] = encodeMentionMessage(mention)
  }
  return msgs
}

func encodeMentionMessage(mention Mention) *pb.Mention {
  return &pb.Mention{
    Indices:     encodeIndicesMessage(mention.Indices),
    UserId:      mention.UserID,
    Handle:      mention.UserHandle,
    DisplayName: mention.UserDisplayName,
  }
}

func decodeMentionMessages(msgs []*pb.Mention) []Mention {
  mentions := make([]Mention, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      mentions[i] = Mention{
        Indices:         decodeIndicesMessage(msg.Indices),
        UserID:          msg.UserId,
        UserHandle:      msg.Handle,
        UserDisplayName: msg.DisplayName,
      }
    }
  }
  return mentions
}

func encodeMediaMessages(medias []Media) []*pb.Media {
  msgs := make([]*pb.Media, len(medias))
  for i, media := range medias {
    msgs[i] = encodeMediaMessage(media)
  }
  return msgs
}

func encodeMediaMessage(media Media) *pb.Media {
  msg := pb.Media{
    Url:      encodeURLMessage(media.URL),
    Id:       media.ID,
    Type:     media.Type,
    MediaUrl: media.MediaURL,
    Alt:      media.Alt,
    Thumb:    encodeMediaSizeMessage(media.Thumb),
    Small:    encodeMediaSizeMessage(media.Small),
    Medium:   encodeMediaSizeMessage(media.Medium),
    Large:    encodeMediaSizeMessage(media.Large),
  }
  if media.SourceTweetID != nil {
    msg.Source = &pb.Media_SourceTweetId{SourceTweetId: *media.SourceTweetID}
  } else {
    msg.Source = &pb.Media_NoSource{}
  }
  return &msg
}

func encodeMediaSizeMessage(mediaSize MediaSize) *pb.Media_Size {
  return &pb.Media_Size{
    Width:  uint32(mediaSize.Width),
    Height: uint32(mediaSize.Height),
    Resize: mediaSize.Resize,
  }
}

func decodeMediaMessages(msgs []*pb.Media) []Media {
  media := make([]Media, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      media[i] = Media{
        URL:           decodeURLMessage(msg.Url),
        ID:            msg.Id,
        Type:          msg.Type,
        MediaURL:      msg.MediaUrl,
        Alt:           msg.Alt,
        SourceTweetID: decodeMediaSourceMessage(msg),
        Thumb:         decodeMediaSizeMessage(msg.Thumb),
        Small:         decodeMediaSizeMessage(msg.Small),
        Medium:        decodeMediaSizeMessage(msg.Medium),
        Large:         decodeMediaSizeMessage(msg.Large),
      }
    }
  }
  return media
}

func decodeMediaSourceMessage(msg *pb.Media) *uint64 {
  if source, ok := msg.Source.(*pb.Media_SourceTweetId); ok && source != nil {
    sourceID := source.SourceTweetId
    return &sourceID
  }
  return nil
}

func decodeMediaSizeMessage(msg *pb.Media_Size) MediaSize {
  if msg == nil {
    return MediaSize{}
  }
  return MediaSize{Width: uint(msg.Width), Height: uint(msg.Height), Resize: msg.Resize}
}

func encodePollMessages(polls []Poll) []*pb.Poll {
  msgs := make([]*pb.Poll, len(polls))
  for i, poll := range polls {
    msgs[i] = encodePollMessage(poll)
  }
  return msgs
}

func encodePollMessage(poll Poll) *pb.Poll {
  msg := pb.Poll{
    DurationMinutes: uint32(math.Floor(poll.Duration.Minutes())),
    EndTime:         poll.EndTime.Unix(),
    Options:         make([]*pb.Poll_Option, len(poll.Options)),
  }
  for i, option := range poll.Options {
    msg.Options[i] = &pb.Poll_Option{
      Position: uint32(option.Position),
      Text:     option.Text,
    }
  }
  return &msg
}

func decodePolls(msgs []*pb.Poll) []Poll {
  polls := make([]Poll, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      polls[i] = Poll{
        EndTime:  time.Unix(msg.EndTime, 0),
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
