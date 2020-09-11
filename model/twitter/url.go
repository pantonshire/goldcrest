package twitter

type URL struct {
  Entity
  URL         string `json:"url"`
  ExpandedURL string `json:"expanded_url"`
  DisplayURL  string `json:"display_url"`
  Status      int    `json:"status"`
  Title       string `json:"title"`
  Description string `json:"description"`
  UnwoundURL  string `json:"unwound_url"`
}
