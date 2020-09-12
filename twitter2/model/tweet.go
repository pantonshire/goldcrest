package model

import "time"

type Tweet struct {
  ID                string    `json:"id"`
  Text              string    `json:"text"`
  AuthorID          string    `json:"author_id"`
  ConversationID    string    `json:"conversation_id"`
  CreatedAt         time.Time `json:"created_at"`
  ReplyUserID       string    `json:"in_reply_to_user_id"`
  Lang              string    `json:"lang"`
  PossiblySensitive bool      `json:"possibly_sensitive"`
  Source            string    `json:"source"`

  Geo struct {
    Coordinates struct {
      Type        string    `json:"type"`
      Coordinates []float64 `json:"coordinates"`
    } `json:"coordinates"`
    PlaceID string `json:"place_id"`
  } `json:"geo"`

  Entities struct {
    Annotations []Annotation `json:"annotations"`
    Cashtags    []Tag        `json:"cashtags"`
    Hashtags    []Tag        `json:"hashtags"`
    Mentions    []Tag        `json:"mentions"`
    URLs        []URL        `json:"urls"`
  } `json:"entities"`

  Attachments []struct {
    PollIDs   []string `json:"poll_ids"`
    MediaKeys []string `json:"media_keys"`
  } `json:"attachments"`

  ContextAnnotations []struct {
    Domain struct {
      ID          string `json:"id"`
      Name        string `json:"name"`
      Description string `json:"description"`
    } `json:"domain"`
    Entity struct {
      ID          string `json:"id"`
      Name        string `json:"name"`
      Description string `json:"description"`
    } `json:"entity"`
  } `json:"context_annotations"`

  PrivateMetrics struct {
    ImpressionCount int `json:"impression_count"`
    LinkClicks      int `json:"url_link_clicks"`
    ProfileClicks   int `json:"user_profile_clicks"`
  } `json:"non_public_metrics"`

  OrganicMetrics struct {
    ImpressionCount int `json:"impression_count"`
    LikeCount       int `json:"like_count"`
    ReplyCount      int `json:"reply_count"`
    RetweetCount    int `json:"retweet_count"`
  } `json:"organic_metrics"`

  PromotedMetrics struct {
    ImpressionCount int `json:"impression_count"`
    LikeCount       int `json:"like_count"`
    ReplyCount      int `json:"reply_count"`
    RetweetCount    int `json:"retweet_count"`
    LinkClicks      int `json:"url_link_clicks"`
    ProfileClicks   int `json:"user_profile_clicks"`
  } `json:"promoted_metrics"`

  PublicMetrics struct {
    RetweetCount int `json:"retweet_count"`
    ReplyCount   int `json:"reply_count"`
    LikeCount    int `json:"like_count"`
    QuoteCount   int `json:"quote_count"`
  } `json:"public_metrics"`

  ReferencedTweets []struct {
    Type string `json:"type"`
    ID   string `json:"id"`
  } `json:"referenced_tweets"`

  Withheld struct {
    Copyright    bool     `json:"copyright"`
    CountryCodes []string `json:"country_codes"`
  } `json:"withheld"`
}
