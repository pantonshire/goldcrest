package model

type SearchResult struct {
  Statuses Timeline `json:"statuses"`
  Meta     struct {
    CompletedIn float32 `json:"completed_in"`
    MaxID       uint64  `json:"max_id"`
    MaxIDStr    string  `json:"max_id_str"`
    NextResults string  `json:"next_results"`
    Query       string  `json:"query"`
    Count       uint    `json:"count"`
    SinceID     uint64  `json:"since_id"`
    SinceIDStr  string  `json:"since_id_str"`
  } `json:"search_metadata"`
}
