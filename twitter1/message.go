package twitter1

import (
  "fmt"
  "github.com/golang/protobuf/ptypes"
  "goldcrest/rpc"
  "goldcrest/twitter1/model"
  "time"
)

func tweetModelToMessage(mod model.Tweet) (*rpc.Tweet, error) {
  var err error
  var msg rpc.Tweet

  msg.Id = mod.ID

  if msg.CreatedAt, err = ptypes.TimestampProto(time.Time(mod.CreatedAt)); err != nil {
    return nil, err
  }

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
  if msg.TextDisplayRange, err = newIndicesMessage(displayTextRange); err != nil {
    return nil, err
  }

  msg.Truncated = mod.Truncated

  msg.Source = mod.Source

  if msg.User, err = userModelToMessage(mod.User); err != nil {
    return nil, err
  }

  if mod.ReplyStatusID != nil && mod.ReplyUserID != nil && mod.ReplyUserScreenName != nil {
    msg.Reply = &rpc.Tweet_RepliedTweet{RepliedTweet: &rpc.Tweet_Reply{
      ReplyToTweetId:    *mod.ReplyStatusID,
      ReplyToUserId:     *mod.ReplyUserID,
      ReplyToUserHandle: *mod.ReplyUserScreenName,
    }}
  } else {
    msg.Reply = &rpc.Tweet_NoReply{NoReply: true}
  }

  if mod.QuotedStatus != nil {
    quote, err := tweetModelToMessage(*mod.QuotedStatus)
    if err != nil {
      return nil, err
    }
    msg.Quote = &rpc.Tweet_QuotedTweet{QuotedTweet: quote}
  } else {
    msg.Quote = &rpc.Tweet_NoQuote{NoQuote: true}
  }

  if mod.RetweetedStatus != nil {
    retweet, err := tweetModelToMessage(*mod.RetweetedStatus)
    if err != nil {
      return nil, err
    }
    msg.Retweet = &rpc.Tweet_RetweetedTweet{RetweetedTweet: retweet}
  } else {
    msg.Retweet = &rpc.Tweet_NoRetweet{NoRetweet: true}
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

  if msg.Hashtags, err = symbolModelsToMessages(entities.Hashtags); err != nil {
    return nil, err
  }

  if msg.Urls, err = urlModelsToMessages(entities.URLs); err != nil {
    return nil, err
  }

  if msg.Mentions, err = mentionModelsToMessages(entities.Mentions); err != nil {
    return nil, err
  }

  if msg.Symbols, err = symbolModelsToMessages(entities.Symbols); err != nil {
    return nil, err
  }

  if msg.Polls, err = pollModelsToMessages(entities.Polls); err != nil {
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
  if msg.Media, err = mediaModelsToMessages(media); err != nil {
    return nil, err
  }

  return &msg, nil
}

func userModelToMessage(mod model.User) (*rpc.User, error) {
  var err error
  var msg rpc.User

  msg.Id = mod.ID

  msg.Handle = mod.ScreenName

  msg.DisplayName = mod.Name

  if msg.CreatedAt, err = ptypes.TimestampProto(time.Time(mod.CreatedAt)); err != nil {
    return nil, err
  }

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

  if msg.UrlUrls, err = urlModelsToMessages(mod.Entities.URL.URLs); err != nil {
    return nil, err
  }

  if msg.BioUrls, err = urlModelsToMessages(mod.Entities.Description.URLs); err != nil {
    return nil, err
  }

  return &msg, nil
}

func newIndicesMessage(indices []uint) (*rpc.Indices, error) {
  if len(indices) != 2 {
    return nil, fmt.Errorf("expected [start,end] index values pair, got %v", indices)
  }
  return &rpc.Indices{Start: uint32(indices[0]), End: uint32(indices[1])}, nil
}

func urlModelToMessage(mod model.URL) (*rpc.URL, error) {
  var err error
  msg := rpc.URL{
    TwitterUrl:  mod.URL,
    DisplayUrl:  mod.DisplayURL,
    ExpandedUrl: mod.ExpandedURL,
  }
  if msg.Indices, err = newIndicesMessage(mod.Indices); err != nil {
    return nil, err
  }
  return &msg, nil
}

func urlModelsToMessages(mods []model.URL) ([]*rpc.URL, error) {
  var err error
  msgs := make([]*rpc.URL, len(mods))
  for i, mod := range mods {
    if msgs[i], err = urlModelToMessage(mod); err != nil {
      return nil, err
    }
  }
  return msgs, nil
}

func symbolModelsToMessages(mods []model.Symbol) ([]*rpc.Symbol, error) {
  var err error
  msgs := make([]*rpc.Symbol, len(mods))
  for i, mod := range mods {
    msg := rpc.Symbol{Text: mod.Text}
    if msg.Indices, err = newIndicesMessage(mod.Indices); err != nil {
      return nil, err
    }
    msgs[i] = &msg
  }
  return msgs, nil
}

func mentionModelsToMessages(mods []model.Mention) ([]*rpc.Mention, error) {
  var err error
  msgs := make([]*rpc.Mention, len(mods))
  for i, mod := range mods {
    msg := rpc.Mention{
      UserId:      mod.ID,
      Handle:      mod.ScreenName,
      DisplayName: mod.Name,
    }
    if msg.Indices, err = newIndicesMessage(mod.Indices); err != nil {
      return nil, err
    }
    msgs[i] = &msg
  }
  return msgs, nil
}

func mediaModelsToMessages(mods []model.Media) ([]*rpc.Media, error) {
  var err error
  msgs := make([]*rpc.Media, len(mods))
  for i, mod := range mods {
    msg := rpc.Media{
      Id:   mod.ID,
      Type: mod.Type,
      Alt:  mod.AltText,
      Thumb: &rpc.Media_MediaSize{
        Width:  uint32(mod.Sizes.Thumb.W),
        Height: uint32(mod.Sizes.Thumb.H),
        Resize: mod.Sizes.Thumb.Resize,
      },
      Small: &rpc.Media_MediaSize{
        Width:  uint32(mod.Sizes.Small.W),
        Height: uint32(mod.Sizes.Small.H),
        Resize: mod.Sizes.Small.Resize,
      },
      Medium: &rpc.Media_MediaSize{
        Width:  uint32(mod.Sizes.Medium.W),
        Height: uint32(mod.Sizes.Medium.H),
        Resize: mod.Sizes.Medium.Resize,
      },
      Large: &rpc.Media_MediaSize{
        Width:  uint32(mod.Sizes.Large.W),
        Height: uint32(mod.Sizes.Large.H),
        Resize: mod.Sizes.Large.Resize,
      },
    }
    if msg.Url, err = urlModelToMessage(mod.URL); err != nil {
      return nil, err
    }
    if mod.MediaURLHttps != "" {
      msg.MediaUrl = mod.MediaURLHttps
    } else {
      msg.MediaUrl = mod.MediaURL
    }
    if mod.SourceStatusID != nil {
      msg.Source = &rpc.Media_SourceTweetId{SourceTweetId: *mod.SourceStatusID}
    } else {
      msg.Source = &rpc.Media_NoSource{NoSource: true}
    }
    msgs[i] = &msg
  }
  return msgs, nil
}

func pollModelsToMessages(mods []model.Poll) ([]*rpc.Poll, error) {
  var err error
  msgs := make([]*rpc.Poll, len(mods))
  for i, mod := range mods {
    msg := rpc.Poll{DurationMinutes: uint32(mod.DurationMinutes)}
    if msg.EndTime, err = ptypes.TimestampProto(time.Time(mod.EndTime)); err != nil {
      return nil, err
    }
    msg.Options = make([]*rpc.Poll_PollOption, len(mod.Options))
    for j, optionMod := range mod.Options {
      msg.Options[j] = &rpc.Poll_PollOption{
        Position: uint32(optionMod.Position),
        Text:     optionMod.Text,
      }
    }
    msgs[i] = &msg
  }
  return msgs, nil
}
