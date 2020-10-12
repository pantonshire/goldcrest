package twitter1

import (
  pb "github.com/pantonshire/goldcrest/proto"
  "github.com/pantonshire/goldcrest/twitter1/model"
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

func decodeTweetModels(ms []model.Tweet) []Tweet {
  tweets := make([]Tweet, len(ms))
  for i, m := range ms {
    tweets[i] = decodeTweetModel(m)
  }
  return tweets
}

func decodeTweetModel(m model.Tweet) Tweet {
  tweet := Tweet{
    ID:                   m.ID,
    CreatedAt:            time.Time(m.CreatedAt),
    Source:               m.Source,
    User:                 decodeUserModel(m.User),
    Replies:              m.ReplyCount,
    Retweets:             m.RetweetCount,
    CurrentUserRetweeted: m.Retweeted,
    FilterLevel:          m.FilterLevel,
    WithheldCopyright:    m.WithheldCopyright,
    WithheldCounties:     m.WithheldInCountries,
  }
  var mEntities model.TweetEntities
  var mExtendedEntities model.TweetExtendedEntities
  if m.ExtendedTweet != nil {
    tweet.Text = m.ExtendedTweet.FullText
    tweet.Truncated = false
    tweet.TextDisplayRange = decodeIndicesModel(m.ExtendedTweet.DisplayTextRange)
    mEntities = m.ExtendedTweet.Entities
    mExtendedEntities = m.ExtendedTweet.ExtendedEntities
  } else {
    if m.FullText != "" {
      tweet.Text = m.FullText
      tweet.Truncated = false
    } else {
      tweet.Text = m.Text
      tweet.Truncated = m.Truncated
    }
    tweet.TextDisplayRange = decodeIndicesModel(m.DisplayTextRange)
    mEntities = m.Entities
    mExtendedEntities = m.ExtendedEntities
  }
  tweet.Hashtags = decodeSymbolModels(mEntities.Hashtags)
  tweet.URLs = decodeURLModels(mEntities.URLs)
  tweet.Mentions = decodeMentionModels(mEntities.Mentions)
  tweet.Symbols = decodeSymbolModels(mEntities.Symbols)
  tweet.Polls = decodePollModels(mEntities.Polls)
  var media []model.Media
  mediaIDs := make(map[uint64]bool)
  for _, mMedia := range mExtendedEntities.Media {
    if !mediaIDs[mMedia.ID] {
      media = append(media, mMedia)
      mediaIDs[mMedia.ID] = true
    }
  }
  for _, mMedia := range mEntities.Media {
    if !mediaIDs[mMedia.ID] {
      media = append(media, mMedia)
      mediaIDs[mMedia.ID] = true
    }
  }
  tweet.Media = decodeMediaModels(media)
  if m.ReplyStatusID != nil && m.ReplyUserID != nil && m.ReplyUserScreenName != nil {
    tweet.RepliedTo = &ReplyData{
      TweetID:    *m.ReplyStatusID,
      UserID:     *m.ReplyUserID,
      UserHandle: *m.ReplyUserScreenName,
    }
  }
  if m.QuotedStatus != nil {
    qt := decodeTweetModel(*m.QuotedStatus)
    tweet.Quoted = &qt
  }
  if m.RetweetedStatus != nil {
    rt := decodeTweetModel(*m.RetweetedStatus)
    tweet.Retweeted = &rt
  }
  if m.QuoteCount != nil {
    tweet.Quotes = *m.QuoteCount
  }
  if m.FavoriteCount != nil {
    tweet.Likes = *m.FavoriteCount
  }
  if m.Favorited != nil {
    tweet.CurrentUserLiked = *m.Favorited
  }
  if m.CurrentUserRetweet != nil {
    rtID := m.CurrentUserRetweet.ID
    tweet.CurrentUserRetweetID = &rtID
  }
  if m.PossiblySensitive != nil {
    tweet.PossiblySensitive = *m.PossiblySensitive
  }
  if m.Lang != nil {
    tweet.Lang = *m.Lang
  }
  if m.WithheldScope != nil {
    tweet.WithheldScope = *m.WithheldScope
  }
  return tweet
}

func encodeTimelineMessage(tweets []Tweet) *pb.Timeline {
  msgs := make([]*pb.Tweet, len(tweets))
  for i, tweet := range tweets {
    msgs[i] = encodeTweetMessage(tweet)
  }
  return &pb.Timeline{Tweets: msgs}
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

func encodeTweetMessage(tweet Tweet) *pb.Tweet {
  msg := pb.Tweet{
    Id:                tweet.ID,
    CreatedAt:         tweet.CreatedAt.Unix(),
    Text:              tweet.Text,
    TextDisplayRange:  encodeIndicesMessage(tweet.TextDisplayRange),
    Truncated:         tweet.Truncated,
    Source:            tweet.Source,
    User:              encodeUserMessage(tweet.User),
    QuoteCount:        uint32(tweet.Quotes),
    ReplyCount:        uint32(tweet.Replies),
    RetweetCount:      uint32(tweet.Retweets),
    FavoriteCount:     uint32(tweet.Likes),
    Favorited:         tweet.CurrentUserLiked,
    Retweeted:         tweet.CurrentUserRetweeted,
    Hashtags:          encodeSymbolMessages(tweet.Hashtags),
    Urls:              encodeURLMessages(tweet.URLs),
    Mentions:          encodeMentionMessages(tweet.Mentions),
    Symbols:           encodeSymbolMessages(tweet.Symbols),
    Media:             encodeMediaMessages(tweet.Media),
    Polls:             encodePollMessages(tweet.Polls),
    PossiblySensitive: tweet.PossiblySensitive,
    FilterLevel:       tweet.FilterLevel,
    Lang:              tweet.Lang,
    WithheldCopyright: tweet.WithheldCopyright,
    WithheldCountries: tweet.WithheldCounties,
    WithheldScope:     tweet.WithheldScope,
  }
  if tweet.RepliedTo != nil {
    msg.Reply = &pb.Tweet_RepliedTweet{RepliedTweet: &pb.Tweet_Reply{
      ReplyToTweetId:    tweet.RepliedTo.TweetID,
      ReplyToUserId:     tweet.RepliedTo.UserID,
      ReplyToUserHandle: tweet.RepliedTo.UserHandle,
    }}
  } else {
    msg.Reply = &pb.Tweet_NoReply{NoReply: true}
  }
  if tweet.Quoted != nil {
    msg.Quote = &pb.Tweet_QuotedTweet{QuotedTweet: encodeTweetMessage(*tweet.Quoted)}
  } else {
    msg.Quote = &pb.Tweet_NoQuote{NoQuote: true}
  }
  if tweet.Retweeted != nil {
    msg.Retweet = &pb.Tweet_RetweetedTweet{RetweetedTweet: encodeTweetMessage(*tweet.Retweeted)}
  } else {
    msg.Retweet = &pb.Tweet_NoRetweet{NoRetweet: true}
  }
  if tweet.CurrentUserRetweetID != nil {
    msg.CurrentUserRetweetId = *tweet.CurrentUserRetweetID
  } else {
    msg.CurrentUserRetweetId = 0
  }
  return &msg
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

func decodeUserModel(m model.User) User {
  user := User{
    ID:                m.ID,
    Handle:            m.ScreenName,
    DisplayName:       m.Name,
    CreatedAt:         time.Time(m.CreatedAt),
    Protected:         m.Protected,
    Verified:          m.Verified,
    FollowerCount:     m.FollowersCount,
    FollowingCount:    m.FriendsCount,
    ListedCount:       m.ListedCount,
    FavouritesCount:   m.FavoritesCount,
    StatusesCount:     m.StatusesCount,
    ProfileBanner:     m.ProfileBanner,
    ProfileImage:      m.ProfileImage,
    DefaultProfile:    m.DefaultProfile,
    WithheldCountries: m.WithheldCountries,
    URLs:              decodeURLModels(m.Entities.URL.URLs),
    BioURLs:           decodeURLModels(m.Entities.Description.URLs),
  }
  if m.Description != nil {
    user.Bio = *m.Description
  }
  if m.URL != nil {
    user.URL = *m.URL
  }
  if m.Location != nil {
    user.Location = *m.Location
  }
  if m.WithheldScope != nil {
    user.WithheldScope = *m.WithheldScope
  }
  return user
}

func encodeUserMessage(user User) *pb.User {
  return &pb.User{
    Id:                  user.ID,
    Handle:              user.Handle,
    DisplayName:         user.DisplayName,
    CreatedAt:           user.CreatedAt.Unix(),
    Bio:                 user.Bio,
    Url:                 user.URL,
    Location:            user.Location,
    Protected:           user.Protected,
    Verified:            user.Verified,
    FollowerCount:       uint32(user.FollowerCount),
    FollowingCount:      uint32(user.FollowerCount),
    ListedCount:         uint32(user.ListedCount),
    FavoritesCount:      uint32(user.FavouritesCount),
    StatusesCount:       uint32(user.StatusesCount),
    ProfileBanner:       user.ProfileBanner,
    ProfileImage:        user.ProfileImage,
    DefaultProfile:      user.DefaultProfile,
    DefaultProfileImage: user.DefaultProfileImage,
    WithheldCountries:   user.WithheldCountries,
    WithheldScope:       user.WithheldScope,
    UrlUrls:             encodeURLMessages(user.URLs),
    BioUrls:             encodeURLMessages(user.BioURLs),
  }
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

func decodeIndicesModel(m model.Indices) Indices {
  var start, end uint
  if len(m) > 0 {
    start = m[0]
    if len(m) > 1 {
      end = m[1]
    }
  }
  return Indices{Start: start, End: end}
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

func decodeURLModels(ms []model.URL) []URL {
  urls := make([]URL, len(ms))
  for i, m := range ms {
    urls[i] = decodeURLModel(m)
  }
  return urls
}

func decodeURLModel(m model.URL) URL {
  return URL{
    Indices:     decodeIndicesModel(m.Indices),
    TwitterURL:  m.URL,
    DisplayURL:  m.DisplayURL,
    ExpandedURL: m.ExpandedURL,
  }
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

func decodeSymbolModels(ms []model.Symbol) []Symbol {
  symbols := make([]Symbol, len(ms))
  for i, m := range ms {
    symbols[i] = decodeSymbolModel(m)
  }
  return symbols
}

func decodeSymbolModel(m model.Symbol) Symbol {
  return Symbol{
    Indices: decodeIndicesModel(m.Indices),
    Text:    m.Text,
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

func decodeMentionModels(ms []model.Mention) []Mention {
  mentions := make([]Mention, len(ms))
  for i, m := range ms {
    mentions[i] = decodeMentionModel(m)
  }
  return mentions
}

func decodeMentionModel(m model.Mention) Mention {
  return Mention{
    Indices:         decodeIndicesModel(m.Indices),
    UserID:          m.ID,
    UserHandle:      m.ScreenName,
    UserDisplayName: m.Name,
  }
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

func decodeMediaModels(ms []model.Media) []Media {
  medias := make([]Media, len(ms))
  for i, m := range ms {
    medias[i] = decodeMediaModel(m)
  }
  return medias
}

func decodeMediaModel(m model.Media) Media {
  media := Media{
    URL:           decodeURLModel(m.URL),
    ID:            m.ID,
    Type:          m.Type,
    Alt:           m.AltText,
    SourceTweetID: m.SourceStatusID,
    Thumb:         decodeMediaSizeModel(m.Sizes.Thumb),
    Small:         decodeMediaSizeModel(m.Sizes.Small),
    Medium:        decodeMediaSizeModel(m.Sizes.Medium),
    Large:         decodeMediaSizeModel(m.Sizes.Large),
  }
  if m.MediaURLHttps != "" {
    media.MediaURL = m.MediaURLHttps
  } else {
    media.MediaURL = m.MediaURL
  }
  return media
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
    msg.Source = &pb.Media_NoSource{NoSource: true}
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

func decodeMediaSizeModel(m model.MediaSize) MediaSize {
  return MediaSize{Width: m.W, Height: m.H, Resize: m.Resize}
}

func decodeMediaSizeMessage(msg *pb.Media_Size) MediaSize {
  if msg == nil {
    return MediaSize{}
  }
  return MediaSize{Width: uint(msg.Width), Height: uint(msg.Height), Resize: msg.Resize}
}

func decodePollModels(ms []model.Poll) []Poll {
  polls := make([]Poll, len(ms))
  for i, m := range ms {
    polls[i] = decodePollModel(m)
  }
  return polls
}

func decodePollModel(m model.Poll) Poll {
  poll := Poll{
    EndTime:  time.Time(m.EndTime),
    Duration: time.Minute * time.Duration(m.DurationMinutes),
    Options:  make([]PollOption, len(m.Options)),
  }
  for i, optM := range m.Options {
    poll.Options[i] = PollOption{
      Position: optM.Position,
      Text:     optM.Text,
    }
  }
  return poll
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
