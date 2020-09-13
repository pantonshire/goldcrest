package model

type Entities struct {
  Hashtags []Hashtag `json:"hashtags"`
  Media    []Media   `json:"media"`
  URLs     []URL     `json:"urls"`
  Mentions []Mention `json:"user_mentions"`
  Symbols  []Symbol  `json:"symbols"`
  Polls    []Poll    `json:"polls"`
}

type ExtendedEntities struct {
  Media []Media `json:"media"`
}

type InlineEntity struct {
  Indices []int `json:"indices"`
}

type Hashtag struct {
  InlineEntity
  Text string `json:"text"`
}

type URL struct {
  InlineEntity
  URL         string `json:"url"`
  DisplayURL  string `json:"display_url"`
  ExpandedURL string `json:"expanded_url"`

  Unwound *struct {
    URL         string `json:"url"`
    Status      int    `json:"status"`
    Title       string `json:"title"`
    Description string `json:"description"`
  } `json:"unwound"`
}

type Mention struct {
  InlineEntity
  ID         int64  `json:"id"`
  IDStr      string `json:"id_str"`
  Name       string `json:"name"`
  ScreenName string `json:"screen_name"`
}

type Symbol struct {
  InlineEntity
  Text string `json:"text"`
}

type Media struct {
  URL
  ID                int64   `json:"id"`
  IDStr             string  `json:"id_str"`
  MediaURL          string  `json:"media_url"`
  MediaURLHttps     string  `json:"media_url_https"`
  Type              string  `json:"type"`
  SourceStatusID    *int64  `json:"source_status_id"`
  SourceStatusIDStr *string `json:"source_status_id_str"`
  AltText           string  `json:"ext_alt_text"`

  Sizes struct {
    Thumb  MediaSize `json:"thumb"`
    Small  MediaSize `json:"small"`
    Medium MediaSize `json:"medium"`
    Large  MediaSize `json:"large"`
  } `json:"sizes"`
}

type MediaSize struct {
  W      int    `json:"w"`
  H      int    `json:"h"`
  Resize string `json:"resize"`
}

type Poll struct {
  EndTime         TwitterTime `json:"end_datetime"`
  DurationMinutes int         `json:"duration_minutes"`

  Options []struct {
    Position int    `json:"position"`
    Text     string `json:"text"`
  } `json:"options"`
}
