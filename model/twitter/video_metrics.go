package twitter

type VideoMetrics struct {
  Playback000Count int `json:"playback_0_count"`
  Playback025Count int `json:"playback_25_count"`
  Playback050Count int `json:"playback_50_count"`
  Playback075Count int `json:"playback_75_count"`
  Playback100Count int `json:"playback_100_count"`
}

type VideoViewMetrics struct {
  VideoMetrics
  ViewCount int `json:"view_count"`
}
