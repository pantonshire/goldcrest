package twitter

type Media struct {
  MediaKey        string           `json:"media_key"`
  Type            string           `json:"type"`
  Width           int              `json:"width"`
  Height          int              `json:"height"`
  DurationMs      int              `json:"duration_ms"`
  PreviewURL      string           `json:"preview_image_url"`
  PrivateMetrics  VideoMetrics     `json:"non_public_metrics"`
  OrganicMetrics  VideoViewMetrics `json:"organic_metrics"`
  PromotedMetrics VideoViewMetrics `json:"promoted_metrics"`
}
