package model

type Timeline []Tweet

type Tweet struct {
  ID    uint64 `json:"id"`
  IDStr string `json:"id_str"`

  CreatedAt TwitterTime `json:"created_at"`

  // The field that will be used to store the tweet text in compatibility mode. The text will be limited to
  // 140 characters, so newer tweets may be cut off; this is indicated by a trailing horizontal ellipsis
  // character (U+2026). This field will be empty if extended mode is being used; the text will appear in
  // FullText in this case. The TextContent function should be preferred over accessing this field directly.
  Text string `json:"text"`

  // The field that will be used to store the tweet text in extended mode. The text will be limited to 280
  // characters. This field will be empty if compatibility mode is being used; the text will appear in Text
  // in this case. The TextContent function should be preferred over accessing this field directly.
  FullText string `json:"full_text"`

  DisplayTextRange Indices `json:"display_text_range"`

  Source    string `json:"source"`
  Truncated bool   `json:"truncated"`

  ReplyStatusID       *uint64 `json:"in_reply_to_status_id"`
  ReplyStatusIDStr    *string `json:"in_reply_to_status_id_str"`
  ReplyUserID         *uint64 `json:"in_reply_to_user_id"`
  ReplyUserIDStr      *string `json:"in_reply_to_user_id_str"`
  ReplyUserScreenName *string `json:"in_reply_to_screen_name"`

  User User `json:"user"`

  Coordinates *struct {
    Coordinates []float64 `json:"coordinates"`
    Type        string    `json:"type"`
  } `json:"coordinates"`

  Place *struct {
    ID          string `json:"id"`
    URL         string `json:"url"`
    PlaceType   string `json:"place_type"`
    Name        string `json:"name"`
    FullName    string `json:"full_name"`
    CountryCode string `json:"country_code"`
    Country     string `json:"country"`

    BoundingBox struct {
      Coordinates [][][]float64 `json:"coordinates"`
      Type        string        `json:"type"`
    } `json:"bounding_box"`
  } `json:"place"`

  QuotedStatusID    *uint64 `json:"quoted_status_id"`
  QuotedStatusIDStr *string `json:"quoted_status_id_str"`
  IsQuoteStatus     bool    `json:"is_quote_status"`

  QuotedStatus    *Tweet `json:"quoted_status"`
  RetweetedStatus *Tweet `json:"retweeted_status"`

  QuoteCount    uint32 `json:"quote_count"`
  ReplyCount    uint32 `json:"reply_count"`
  RetweetCount  uint32 `json:"retweet_count"`
  FavoriteCount uint32 `json:"favorite_count"`

  Entities         TweetEntities         `json:"entities"`
  ExtendedEntities TweetExtendedEntities `json:"extended_entities"`

  Favorited         bool   `json:"favorited"`
  Retweeted         bool   `json:"retweeted"`
  PossiblySensitive bool   `json:"possibly_sensitive"`
  FilterLevel       string `json:"filter_level"`
  Lang              string `json:"lang"`

  WithheldCopyright   bool     `json:"withheld_copyright"`
  WithheldInCountries []string `json:"withheld_in_countries"`
  WithheldScope       string   `json:"withheld_scope"`

  MatchingRules []struct {
    Tag   string `json:"tag"`
    ID    uint64 `json:"id"`
    IDStr string `json:"id_str"`
  } `json:"matching_rules"`

  CurrentUserRetweet *struct {
    ID    uint64 `json:"id"`
    IDStr string `json:"id_str"`
  } `json:"current_user_retweet"`

  Scopes *struct {
    Followers bool `json:"followers"`
  } `json:"scopes"`

  // Included in streamed tweets
  ExtendedTweet *ExtendedTweet `json:"extended_tweet"`
}

type ExtendedTweet struct {
  FullText         string                `json:"full_text"`
  DisplayTextRange Indices               `json:"display_text_range"`
  Entities         TweetEntities         `json:"entities"`
  ExtendedEntities TweetExtendedEntities `json:"extended_entities"`
}

// Returns the text of the Tweet. Since the text may be stored in one of a number of different fields, this
// method should be preferred over directly accessing fields in Tweet relating to text.
func (tweet Tweet) TextContent() (string, Indices) {
  if tweet.ExtendedTweet != nil && tweet.ExtendedTweet.FullText != "" {
    return tweet.ExtendedTweet.FullText, tweet.ExtendedTweet.DisplayTextRange
  } else if tweet.FullText != "" {
    return tweet.FullText, tweet.DisplayTextRange
  }
  return tweet.Text, tweet.DisplayTextRange
}

func (tweet Tweet) EntityContent() TweetEntities {
  var entities TweetEntities
  var extendedEntities TweetExtendedEntities
  if tweet.ExtendedTweet != nil {
    entities = tweet.ExtendedTweet.Entities
    extendedEntities = tweet.ExtendedTweet.ExtendedEntities
  } else {
    entities = tweet.Entities
    extendedEntities = tweet.ExtendedEntities
  }
updateMedia:
  for _, extendedMedia := range extendedEntities.Media {
    for i := 0; i < len(entities.Media); i++ {
      if extendedMedia.ID == entities.Media[i].ID {
        entities.Media[i] = extendedMedia
        continue updateMedia
      }
    }
    entities.Media = append(entities.Media, extendedMedia)
  }
  return entities
}
