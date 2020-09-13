package model

type Tweet struct {
  ID    int64  `json:"id"`
  IDStr string `json:"id_str"`

  CreatedAt TwitterTime `json:"created_at"`

  //Compatibility mode text, 140 character limit
  Text string `json:"text"`

  //Extended mode text, 280 character limit
  FullText string `json:"full_text"`

  DisplayTextRange []int `json:"display_text_range"`

  Source    string `json:"source"`
  Truncated bool   `json:"truncated"`

  ReplyStatusID       *int64  `json:"in_reply_to_status_id"`
  ReplyStatusIDStr    *string `json:"in_reply_to_status_id_str"`
  ReplyUserID         *int64  `json:"in_reply_to_user_id"`
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

  QuotedStatusID    *int64  `json:"quoted_status_id"`
  QuotedStatusIDStr *string `json:"quoted_status_id_str"`
  IsQuoteStatus     bool    `json:"is_quote_status"`

  QuotedStatus    *Tweet `json:"quoted_status"`
  RetweetedStatus *Tweet `json:"retweeted_status"`

  QuoteCount    *int `json:"quote_count"`
  ReplyCount    int  `json:"reply_count"`
  RetweetCount  int  `json:"retweet_count"`
  FavoriteCount *int `json:"favorite_count"`

  Entities struct {
    Hashtags []Hashtag `json:"hashtags"`
    Media    []Media   `json:"media"`
    URLs     []URL     `json:"urls"`
    Mentions []Mention `json:"user_mentions"`
    Symbols  []Symbol  `json:"symbols"`
    Polls    []Poll    `json:"polls"`
  } `json:"entities"`

  ExtendedEntities struct {
    Media []Media `json:"media"`
  } `json:"extended_entities"`

  Favorited         *bool   `json:"favorited"`
  Retweeted         bool    `json:"retweeted"`
  PossiblySensitive *bool   `json:"possibly_sensitive"`
  FilterLevel       string  `json:"filter_level"`
  Lang              *string `json:"lang"`

  WithheldCopyright   bool     `json:"withheld_copyright"`
  WithheldInCountries []string `json:"withheld_in_countries"`
  WithheldScope       *string  `json:"withheld_scope"`

  MatchingRules []struct {
    Tag   string `json:"tag"`
    ID    int64  `json:"id"`
    IDStr string `json:"id_str"`
  } `json:"matching_rules"`

  CurrentUserRetweet *struct {
    ID    uint64 `json:"id"`
    IDStr string `json:"id_str"`
  } `json:"current_user_retweet"`

  Scopes *struct {
    Followers bool `json:"followers"`
  } `json:"scopes"`

  //Included in streamed tweets
  ExtendedTweet *struct {
    FullText         string           `json:"full_text"`
    DisplayTextRange []int            `json:"display_text_range"`
    Entities         Entities         `json:"entities"`
    ExtendedEntities ExtendedEntities `json:"extended_entities"`
  } `json:"extended_tweet"`
}
