package model

// The result of media/upload INIT
type MediaInit struct {
  MediaID     uint64 `json:"media_id"`
  MediaIDStr  string `json:"media_id_string"`
  Size        uint   `json:"size"`
  ExpiresSecs uint   `json:"expires_after_secs"`
  Image       struct {
    ImageType string `json:"image_type"`
    W         uint   `json:"w"`
    H         uint   `json:"h"`
  } `json:"image"`
}

// The result of media/upload STATUS and FINALIZE
type MediaStatus struct {
  MediaID        uint64 `json:"media_id"`
  MediaIDStr     string `json:"media_id_string"`
  ExpiresSecs    uint   `json:"expires_after_secs"`
  ProcessingInfo *struct { // absent from FINALIZE response if no STATUS calls needed
    State           string `json:"state"`
    CheckAfterSecs  uint   `json:"check_after_secs"` // how long to wait before checking STATUS
    ProgressPercent uint   `json:"progress_percent"`
    Error           *struct {
      Code    int    `json:"code"`
      Name    string `json:"name"`
      Message string `json:"message"`
    } `json:"error"`
  } `json:"processing_info"`
}
