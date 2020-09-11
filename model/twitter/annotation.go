package twitter

type Annotation struct {
  Entity
  Probability    float64 `json:"probability"`
  Type           string  `json:"type"`
  NormalizedText string  `json:"normalized_text"`
}
