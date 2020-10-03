package twitter1

import (
  "fmt"
  pb "goldcrest/proto"
  "goldcrest/twitter1/model"
  "time"
)

func encodeAuthPair(auth AuthPair) *pb.Authentication {
  return &pb.Authentication{
    ConsumerKey: auth.Public.Key,
    AccessToken: auth.Public.Token,
    SecretKey:   auth.Secret.Key,
    SecretToken: auth.Secret.Token,
  }
}

func decodeAuthPair(authMessage *pb.Authentication) AuthPair {
  if authMessage == nil {
    return AuthPair{}
  }
  return AuthPair{
    Auth{Key: authMessage.SecretKey, Token: authMessage.SecretToken},
    Auth{Key: authMessage.ConsumerKey, Token: authMessage.AccessToken},
  }
}

func encodeTweetOptions(params TweetOptions) *pb.TweetOptions {
  return &pb.TweetOptions{
    TrimUser:          params.TrimUser,
    IncludeMyRetweet:  params.IncludeMyRetweet,
    IncludeEntities:   params.IncludeEntities,
    IncludeExtAltText: params.IncludeExtAltText,
    IncludeCardUri:    params.TrimUser,
    Mode:              encodeTweetMode(params.Mode),
  }
}

func decodeTweetOptions(optsMessage *pb.TweetOptions) TweetOptions {
  if optsMessage == nil {
    return TweetOptions{}
  }
  return TweetOptions{
    TrimUser:          optsMessage.TrimUser,
    IncludeMyRetweet:  optsMessage.IncludeMyRetweet,
    IncludeEntities:   optsMessage.IncludeEntities,
    IncludeExtAltText: optsMessage.IncludeExtAltText,
    IncludeCardURI:    optsMessage.TrimUser,
    Mode:              decodeTweetMode(optsMessage.Mode),
  }
}

func encodeTweetMode(mode TweetMode) pb.TweetOptions_Mode {
  if mode == ExtendedMode {
    return pb.TweetOptions_EXTENDED
  }
  return pb.TweetOptions_COMPAT
}

func decodeTweetMode(mode pb.TweetOptions_Mode) TweetMode {
  if mode == pb.TweetOptions_EXTENDED {
    return ExtendedMode
  }
  return CompatibilityMode
}

func encodeTweets(mods []model.Tweet) (*pb.Timeline, error) {
  msgs := make([]*pb.Tweet, len(mods))
  for i, mod := range mods {
    msg, err := encodeTweet(mod)
    if err != nil {
      return nil, err
    }
    msgs[i] = msg
  }
  return &pb.Timeline{Tweets: msgs}, nil
}

func decodeTimeline(msg *pb.Timeline) []Tweet {
  if msg == nil {
    return nil
  }
  tweets := make([]Tweet, len(msg.Tweets))
  for i, tweetMsg := range msg.Tweets {
    tweets[i] = decodeTweet(tweetMsg)
  }
  return tweets
}

func encodeTweet(mod model.Tweet) (*pb.Tweet, error) {
  var err error
  var msg pb.Tweet

  msg.Id = mod.ID

  msg.CreatedAt = time.Time(mod.CreatedAt).Unix()

  if mod.ExtendedTweet != nil {
    msg.Text = mod.ExtendedTweet.FullText
  } else if mod.FullText != "" {
    msg.Text = mod.FullText
  } else {
    msg.Text = mod.Text
  }

  var displayTextRange []uint
  if mod.ExtendedTweet != nil {
    displayTextRange = mod.ExtendedTweet.DisplayTextRange
  } else {
    displayTextRange = mod.DisplayTextRange
  }
  if msg.TextDisplayRange, err = encodeIndices(displayTextRange); err != nil {
    return nil, err
  }

  msg.Truncated = mod.Truncated

  msg.Source = mod.Source

  if msg.User, err = encodeUser(mod.User); err != nil {
    return nil, err
  }

  if mod.ReplyStatusID != nil && mod.ReplyUserID != nil && mod.ReplyUserScreenName != nil {
    msg.Reply = &pb.Tweet_RepliedTweet{RepliedTweet: &pb.Tweet_Reply{
      ReplyToTweetId:    *mod.ReplyStatusID,
      ReplyToUserId:     *mod.ReplyUserID,
      ReplyToUserHandle: *mod.ReplyUserScreenName,
    }}
  } else {
    msg.Reply = &pb.Tweet_NoReply{NoReply: true}
  }

  if mod.QuotedStatus != nil {
    quote, err := encodeTweet(*mod.QuotedStatus)
    if err != nil {
      return nil, err
    }
    msg.Quote = &pb.Tweet_QuotedTweet{QuotedTweet: quote}
  } else {
    msg.Quote = &pb.Tweet_NoQuote{NoQuote: true}
  }

  if mod.RetweetedStatus != nil {
    retweet, err := encodeTweet(*mod.RetweetedStatus)
    if err != nil {
      return nil, err
    }
    msg.Retweet = &pb.Tweet_RetweetedTweet{RetweetedTweet: retweet}
  } else {
    msg.Retweet = &pb.Tweet_NoRetweet{NoRetweet: true}
  }

  if mod.QuoteCount != nil {
    msg.QuoteCount = uint32(*mod.QuoteCount)
  }

  msg.ReplyCount = uint32(mod.ReplyCount)

  msg.RetweetCount = uint32(mod.RetweetCount)

  if mod.FavoriteCount != nil {
    msg.FavoriteCount = uint32(*mod.FavoriteCount)
  }

  if mod.CurrentUserRetweet != nil {
    msg.CurrentUserRetweetId = mod.CurrentUserRetweet.ID
  }

  var entities model.TweetEntities
  var extendedEntities model.TweetExtendedEntities

  if mod.ExtendedTweet != nil {
    entities, extendedEntities = mod.ExtendedTweet.Entities, mod.ExtendedTweet.ExtendedEntities
  } else {
    entities, extendedEntities = mod.Entities, mod.ExtendedEntities
  }

  if msg.Hashtags, err = encodeSymbols(entities.Hashtags); err != nil {
    return nil, err
  }

  if msg.Urls, err = encodeURLs(entities.URLs); err != nil {
    return nil, err
  }

  if msg.Mentions, err = encodeMentions(entities.Mentions); err != nil {
    return nil, err
  }

  if msg.Symbols, err = encodeSymbols(entities.Symbols); err != nil {
    return nil, err
  }

  if msg.Polls, err = encodePolls(entities.Polls); err != nil {
    return nil, err
  }

  var media []model.Media
  mediaIDs := make(map[uint64]bool)
  for _, mm := range extendedEntities.Media {
    if !mediaIDs[mm.ID] {
      media = append(media, mm)
      mediaIDs[mm.ID] = true
    }
  }
  for _, mm := range entities.Media {
    if !mediaIDs[mm.ID] {
      media = append(media, mm)
      mediaIDs[mm.ID] = true
    }
  }
  if msg.Media, err = encodeMedia(media); err != nil {
    return nil, err
  }

  return &msg, nil
}

func decodeTweet(msg *pb.Tweet) Tweet {
  if msg == nil {
    return Tweet{}
  }
  tweet := Tweet{
    ID:                   msg.Id,
    CreatedAt:            time.Unix(msg.CreatedAt, 0),
    Text:                 msg.Text,
    TextDisplayRange:     decodeIndices(msg.TextDisplayRange),
    Truncated:            msg.Truncated,
    Source:               msg.Source,
    User:                 decodeUser(msg.User),
    Quotes:               uint(msg.QuoteCount),
    Replies:              uint(msg.ReplyCount),
    Retweets:             uint(msg.RetweetCount),
    Likes:                uint(msg.FavoriteCount),
    CurrentUserLiked:     msg.Favorited,
    CurrentUserRetweeted: msg.Retweeted,
    Hashtags:             decodeSymbols(msg.Hashtags),
    URLs:                 decodeURLs(msg.Urls),
    Mentions:             decodeMentions(msg.Mentions),
    Symbols:              decodeSymbols(msg.Symbols),
    Media:                decodeMedia(msg.Media),
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
    decodedQuote := decodeTweet(quote.QuotedTweet)
    tweet.Quoted = &decodedQuote
  }
  if retweet, ok := msg.Retweet.(*pb.Tweet_RetweetedTweet); ok && retweet != nil {
    decodedRetweet := decodeTweet(retweet.RetweetedTweet)
    tweet.Retweeted = &decodedRetweet
  }
  if msg.CurrentUserRetweetId != 0 {
    retweetID := msg.CurrentUserRetweetId
    tweet.CurrentUserRetweetID = &retweetID
  }
  return tweet
}

func encodeUser(mod model.User) (*pb.User, error) {
  var err error
  var msg pb.User

  msg.Id = mod.ID

  msg.Handle = mod.ScreenName

  msg.DisplayName = mod.Name

  msg.CreatedAt = time.Time(mod.CreatedAt).Unix()

  if mod.Description != nil {
    msg.Bio = *mod.Description
  }

  if mod.URL != nil {
    msg.Url = *mod.URL
  }

  if mod.Location != nil {
    msg.Location = *mod.Location
  }

  msg.Protected = mod.Protected

  msg.Verified = mod.Verified

  msg.FollowerCount = uint32(mod.FollowersCount)

  msg.FollowingCount = uint32(mod.FriendsCount)

  msg.ListedCount = uint32(mod.ListedCount)

  msg.FavoritesCount = uint32(mod.FavoritesCount)

  msg.StatusesCount = uint32(mod.StatusesCount)

  msg.ProfileBanner = mod.ProfileBanner

  msg.ProfileImage = mod.ProfileImage

  msg.DefaultProfile = mod.DefaultProfile

  msg.DefaultProfileImage = mod.DefaultProfileImage

  msg.WithheldCountries = mod.WithheldCountries

  if mod.WithheldScope != nil {
    msg.WithheldScope = *mod.WithheldScope
  }

  if msg.UrlUrls, err = encodeURLs(mod.Entities.URL.URLs); err != nil {
    return nil, err
  }

  if msg.BioUrls, err = encodeURLs(mod.Entities.Description.URLs); err != nil {
    return nil, err
  }

  return &msg, nil
}

func decodeUser(msg *pb.User) User {
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
    URLs:                decodeURLs(msg.UrlUrls),
    BioURLs:             decodeURLs(msg.BioUrls),
  }
}

func encodeIndices(indices []uint) (*pb.Indices, error) {
  if len(indices) == 0 {
    return &pb.Indices{Start: 0, End: 0}, nil
  } else if len(indices) != 2 {
    return nil, fmt.Errorf("expected [start,end] index values pair, got %v", indices)
  }
  return &pb.Indices{Start: uint32(indices[0]), End: uint32(indices[1])}, nil
}

func decodeIndices(msg *pb.Indices) Indices {
  if msg == nil {
    return Indices{}
  }
  return Indices{Start: uint(msg.Start), End: uint(msg.End)}
}

func encodeURL(mod model.URL) (*pb.URL, error) {
  var err error
  msg := pb.URL{
    TwitterUrl:  mod.URL,
    DisplayUrl:  mod.DisplayURL,
    ExpandedUrl: mod.ExpandedURL,
  }
  if msg.Indices, err = encodeIndices(mod.Indices); err != nil {
    return nil, err
  }
  return &msg, nil
}

func encodeURLs(mods []model.URL) ([]*pb.URL, error) {
  var err error
  msgs := make([]*pb.URL, len(mods))
  for i, mod := range mods {
    if msgs[i], err = encodeURL(mod); err != nil {
      return nil, err
    }
  }
  return msgs, nil
}

func decodeURL(msg *pb.URL) URL {
  if msg == nil {
    return URL{}
  }
  return URL{
    Indices:     decodeIndices(msg.Indices),
    TwitterURL:  msg.TwitterUrl,
    DisplayURL:  msg.DisplayUrl,
    ExpandedURL: msg.ExpandedUrl,
  }
}

func decodeURLs(msgs []*pb.URL) []URL {
  urls := make([]URL, len(msgs))
  for i, msg := range msgs {
    urls[i] = decodeURL(msg)
  }
  return urls
}

func encodeSymbols(mods []model.Symbol) ([]*pb.Symbol, error) {
  var err error
  msgs := make([]*pb.Symbol, len(mods))
  for i, mod := range mods {
    msg := pb.Symbol{Text: mod.Text}
    if msg.Indices, err = encodeIndices(mod.Indices); err != nil {
      return nil, err
    }
    msgs[i] = &msg
  }
  return msgs, nil
}

func decodeSymbols(msgs []*pb.Symbol) []Symbol {
  symbols := make([]Symbol, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      symbols[i] = Symbol{
        Indices: decodeIndices(msg.Indices),
        Text:    msg.Text,
      }
    }
  }
  return symbols
}

func encodeMentions(mods []model.Mention) ([]*pb.Mention, error) {
  var err error
  msgs := make([]*pb.Mention, len(mods))
  for i, mod := range mods {
    msg := pb.Mention{
      UserId:      mod.ID,
      Handle:      mod.ScreenName,
      DisplayName: mod.Name,
    }
    if msg.Indices, err = encodeIndices(mod.Indices); err != nil {
      return nil, err
    }
    msgs[i] = &msg
  }
  return msgs, nil
}

func decodeMentions(msgs []*pb.Mention) []Mention {
  mentions := make([]Mention, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      mentions[i] = Mention{
        Indices:         decodeIndices(msg.Indices),
        UserID:          msg.UserId,
        UserHandle:      msg.Handle,
        UserDisplayName: msg.DisplayName,
      }
    }
  }
  return mentions
}

func encodeMedia(mods []model.Media) ([]*pb.Media, error) {
  var err error
  msgs := make([]*pb.Media, len(mods))
  for i, mod := range mods {
    msg := pb.Media{
      Id:   mod.ID,
      Type: mod.Type,
      Alt:  mod.AltText,
      Thumb: &pb.Media_Size{
        Width:  uint32(mod.Sizes.Thumb.W),
        Height: uint32(mod.Sizes.Thumb.H),
        Resize: mod.Sizes.Thumb.Resize,
      },
      Small: &pb.Media_Size{
        Width:  uint32(mod.Sizes.Small.W),
        Height: uint32(mod.Sizes.Small.H),
        Resize: mod.Sizes.Small.Resize,
      },
      Medium: &pb.Media_Size{
        Width:  uint32(mod.Sizes.Medium.W),
        Height: uint32(mod.Sizes.Medium.H),
        Resize: mod.Sizes.Medium.Resize,
      },
      Large: &pb.Media_Size{
        Width:  uint32(mod.Sizes.Large.W),
        Height: uint32(mod.Sizes.Large.H),
        Resize: mod.Sizes.Large.Resize,
      },
    }
    if msg.Url, err = encodeURL(mod.URL); err != nil {
      return nil, err
    }
    if mod.MediaURLHttps != "" {
      msg.MediaUrl = mod.MediaURLHttps
    } else {
      msg.MediaUrl = mod.MediaURL
    }
    if mod.SourceStatusID != nil {
      msg.Source = &pb.Media_SourceTweetId{SourceTweetId: *mod.SourceStatusID}
    } else {
      msg.Source = &pb.Media_NoSource{NoSource: true}
    }
    msgs[i] = &msg
  }
  return msgs, nil
}

func decodeMedia(msgs []*pb.Media) []Media {
  media := make([]Media, len(msgs))
  for i, msg := range msgs {
    if msg != nil {
      media[i] = Media{
        URL:           decodeURL(msg.Url),
        ID:            msg.Id,
        Type:          msg.Type,
        MediaURL:      msg.MediaUrl,
        Alt:           msg.Alt,
        SourceTweetID: decodeMediaSource(msg),
        Thumb:         decodeMediaSize(msg.Thumb),
        Small:         decodeMediaSize(msg.Small),
        Medium:        decodeMediaSize(msg.Medium),
        Large:         decodeMediaSize(msg.Large),
      }
    }
  }
  return media
}

func decodeMediaSource(msg *pb.Media) *uint64 {
  if source, ok := msg.Source.(*pb.Media_SourceTweetId); ok && source != nil {
    sourceID := source.SourceTweetId
    return &sourceID
  }
  return nil
}

func decodeMediaSize(msg *pb.Media_Size) MediaSize {
  if msg == nil {
    return MediaSize{}
  }
  return MediaSize{Width: uint(msg.Width), Height: uint(msg.Height), Resize: msg.Resize}
}

func encodePolls(mods []model.Poll) ([]*pb.Poll, error) {
  msgs := make([]*pb.Poll, len(mods))
  for i, mod := range mods {
    msg := pb.Poll{
      DurationMinutes: uint32(mod.DurationMinutes),
      EndTime:         time.Time(mod.EndTime).Unix(),
    }
    msg.Options = make([]*pb.Poll_Option, len(mod.Options))
    for j, optionMod := range mod.Options {
      msg.Options[j] = &pb.Poll_Option{
        Position: uint32(optionMod.Position),
        Text:     optionMod.Text,
      }
    }
    msgs[i] = &msg
  }
  return msgs, nil
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
