package proxy

import (
  pb "github.com/pantonshire/goldcrest/protocol"
  "github.com/pantonshire/goldcrest/proxy/model"
)

func serTimeline(mods model.Timeline) *pb.Timeline {
  tweets := make([]*pb.Tweet, len(mods))
  for i, mod := range mods {
    tweets[i] = serTweet(mod)
  }
  return &pb.Timeline{Tweets: tweets}
}

func serTweet(mod model.Tweet) *pb.Tweet {
  text, displayRange := mod.TextContent()
  entities := mod.EntityContent()
  msg := pb.Tweet{
    Id:                mod.ID,
    CreatedAt:         mod.CreatedAt.Unix(),
    Text:              text,
    TextDisplayRange:  serIndices(displayRange),
    Truncated:         mod.Truncated,
    Source:            mod.Source,
    User:              serUser(mod.User),
    QuoteCount:        mod.QuoteCount,
    ReplyCount:        mod.ReplyCount,
    RetweetCount:      mod.RetweetCount,
    FavoriteCount:     mod.FavoriteCount,
    Favorited:         mod.Favorited,
    Retweeted:         mod.Retweeted,
    Hashtags:          serSymbols(entities.Hashtags),
    Urls:              serURLs(entities.URLs),
    Mentions:          serMentions(entities.Mentions),
    Symbols:           serSymbols(entities.Symbols),
    Media:             serMediaItems(entities.Media),
    Polls:             serPolls(entities.Polls),
    PossiblySensitive: mod.PossiblySensitive,
    FilterLevel:       mod.FilterLevel,
    Lang:              mod.Lang,
    WithheldCopyright: mod.WithheldCopyright,
    WithheldCountries: mod.WithheldInCountries,
    WithheldScope:     mod.WithheldScope,
  }
  if mod.ReplyStatusID != nil && mod.ReplyUserID != nil {
    msg.Reply = &pb.Tweet_RepliedTweet{RepliedTweet: &pb.Tweet_Reply{
      ReplyToTweetId:    *mod.ReplyStatusID,
      ReplyToUserId:     *mod.ReplyUserID,
      ReplyToUserHandle: strSafeDeref(mod.ReplyUserScreenName),
    }}
  } else {
    msg.Reply = &pb.Tweet_NoReply{}
  }
  if mod.QuotedStatus != nil {
    msg.Quote = &pb.Tweet_QuotedTweet{QuotedTweet: serTweet(*mod.QuotedStatus)}
  } else {
    msg.Quote = &pb.Tweet_NoQuote{}
  }
  if mod.RetweetedStatus != nil {
    msg.Retweet = &pb.Tweet_RetweetedTweet{RetweetedTweet: serTweet(*mod.RetweetedStatus)}
  } else {
    msg.Retweet = &pb.Tweet_NoRetweet{}
  }
  if mod.CurrentUserRetweet != nil {
    msg.CurrentUserRetweet = &pb.Tweet_CurrentUserRetweetId{CurrentUserRetweetId: mod.CurrentUserRetweet.ID}
  } else {
    msg.CurrentUserRetweet = &pb.Tweet_NoCurrentUserRetweet{}
  }
  return &msg
}

func serUser(mod model.User) *pb.User {
  return &pb.User{
    Id:                  mod.ID,
    Handle:              mod.ScreenName,
    DisplayName:         mod.Name,
    CreatedAt:           mod.CreatedAt.Unix(),
    Bio:                 mod.Description,
    Url:                 mod.URL,
    Location:            mod.Location,
    Protected:           mod.Protected,
    Verified:            mod.Verified,
    FollowerCount:       mod.FollowersCount,
    FollowingCount:      mod.FriendsCount,
    ListedCount:         mod.ListedCount,
    FavoritesCount:      mod.FavoritesCount,
    StatusesCount:       mod.StatusesCount,
    ProfileBanner:       mod.ProfileBanner,
    ProfileImage:        mod.ProfileImage,
    DefaultProfile:      mod.DefaultProfile,
    DefaultProfileImage: mod.DefaultProfileImage,
    WithheldCountries:   mod.WithheldCountries,
    WithheldScope:       mod.WithheldScope,
    UrlUrls:             serURLs(mod.Entities.URL.URLs),
    BioUrls:             serURLs(mod.Entities.Description.URLs),
  }
}

func serIndices(mod model.Indices) *pb.Indices {
  var start, end uint32
  if len(mod) > 0 {
    start = mod[0]
    if len(mod) > 1 {
      end = mod[1]
    }
  }
  return &pb.Indices{
    Start: start,
    End:   end,
  }
}

func serURLs(mods []model.URL) []*pb.URL {
  msgs := make([]*pb.URL, len(mods))
  for i, mod := range mods {
    msgs[i] = serURL(mod)
  }
  return msgs
}

func serURL(mod model.URL) *pb.URL {
  return &pb.URL{
    Indices:     serIndices(mod.Indices),
    TwitterUrl:  mod.URL,
    DisplayUrl:  mod.DisplayURL,
    ExpandedUrl: mod.ExpandedURL,
  }
}

func serSymbols(mods []model.Symbol) []*pb.Symbol {
  msgs := make([]*pb.Symbol, len(mods))
  for i, mod := range mods {
    msgs[i] = serSymbol(mod)
  }
  return msgs
}

func serSymbol(mod model.Symbol) *pb.Symbol {
  return &pb.Symbol{
    Indices: serIndices(mod.Indices),
    Text:    mod.Text,
  }
}

func serMentions(mods []model.Mention) []*pb.Mention {
  msgs := make([]*pb.Mention, len(mods))
  for i, mod := range mods {
    msgs[i] = serMention(mod)
  }
  return msgs
}

func serMention(mod model.Mention) *pb.Mention {
  return &pb.Mention{
    Indices:     serIndices(mod.Indices),
    UserId:      mod.ID,
    Handle:      mod.ScreenName,
    DisplayName: mod.Name,
  }
}

func serMediaItems(mods []model.Media) []*pb.Media {
  msgs := make([]*pb.Media, len(mods))
  for i, mod := range mods {
    msgs[i] = serMediaItem(mod)
  }
  return msgs
}

func serMediaItem(mod model.Media) *pb.Media {
  msg := pb.Media{
    Url:      serURL(mod.URL),
    Id:       mod.ID,
    Type:     mod.Type,
    MediaUrl: strAlt(mod.MediaURLHttps, mod.MediaURL),
    Alt:      mod.AltText,
    Thumb:    serMediaSize(mod.Sizes.Thumb),
    Small:    serMediaSize(mod.Sizes.Small),
    Medium:   serMediaSize(mod.Sizes.Medium),
    Large:    serMediaSize(mod.Sizes.Large),
  }
  if mod.SourceStatusID != nil {
    msg.Source = &pb.Media_SourceTweetId{SourceTweetId: *mod.SourceStatusID}
  } else {
    msg.Source = &pb.Media_NoSource{}
  }
  return &msg
}

func serMediaSize(mod model.MediaSize) *pb.Media_Size {
  return &pb.Media_Size{
    Width:  mod.W,
    Height: mod.H,
    Resize: mod.Resize,
  }
}

func serPolls(mods []model.Poll) []*pb.Poll {
  msgs := make([]*pb.Poll, len(mods))
  for i, mod := range mods {
    msgs[i] = serPoll(mod)
  }
  return msgs
}

func serPoll(mod model.Poll) *pb.Poll {
  return &pb.Poll{
    DurationMinutes: mod.DurationMinutes,
    EndTime:         mod.EndTime.Unix(),
    Options:         serPollOptions(mod.Options),
  }
}

func serPollOptions(mods []model.PollOption) []*pb.Poll_Option {
  msgs := make([]*pb.Poll_Option, len(mods))
  for i, mod := range mods {
    msgs[i] = &pb.Poll_Option{
      Position: mod.Position,
      Text:     mod.Text,
    }
  }
  return msgs
}
